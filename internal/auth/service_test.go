// internal/auth/service_test.go
package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//
//  TEST SETUP: FAKES
//
// Small fake implementations of UserRepository and PasswordHasher are used so that:
//
// - SignUp can run without a real database.
// - Tests can inspect what the service tried to do (for example, which user
//   it created and what password was stored).
// - Tests stay as pure unit tests of the service logic.
//

// fakeUserRepo records any users that the service attempts to create.
// No real DB access is performed.
type fakeUserRepo struct {
	createdUsers []User
}

// Create appends the user to createdUsers so tests can inspect the result.
func (f *fakeUserRepo) Create(ctx context.Context, u User) error {
	f.createdUsers = append(f.createdUsers, u)
	return nil
}

// FindByEmail is not used in these tests, but is implemented to satisfy
// the UserRepository interface. For now it always returns "not found".
func (f *fakeUserRepo) FindByEmail(ctx context.Context, email string) (User, error) {
	return User{}, ErrUserNotFound
}

// fakeHasher simulates the password hasher.
//
// Instead of running a real hash (such as bcrypt), it simply prepends "HASHED:"
// so tests can easily verify that:
// - The service is not storing the raw password.
// - The service is actually calling the hasher.
type fakeHasher struct{}

// HashPassword returns a deterministic "hashed" version of the password.
func (f *fakeHasher) HashPassword(pw string) (string, error) {
	return "HASHED:" + pw, nil
}

// CheckPassword returns nil if the "hash" matches the expected format.
func (f *fakeHasher) CheckPassword(hash, pw string) error {
	if hash == "HASHED:"+pw {
		return nil
	}
	return ErrInvalidCredentials
}

// newTestAuthService builds an authService using in-memory fakes.
// This is the default constructor for tests that do not need to inspect
// repository calls.
func newTestAuthService() *authService {
	return NewAuthService(&fakeUserRepo{}, &fakeHasher{}).(*authService)
}

// newTestAuthServiceWithFakeRepo builds an authService using a specific
// fakeUserRepo, allowing tests to inspect createdUsers.
func newTestAuthServiceWithFakeRepo(repo *fakeUserRepo) *authService {
	return NewAuthService(repo, &fakeHasher{}).(*authService)
}

//
//  TESTS
//

// TestSignUp_InvalidEmail
// -----------------------
// Given an invalid email string, the service should reject it and return
// ErrInvalidEmailFormat before any call to the repository.
func TestSignUp_InvalidEmail(t *testing.T) {
	svc := newTestAuthService()

	err := svc.SignUp(context.Background(), "not-an-email", "Password123")

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidEmailFormat)
}

// TestSignUp_WeakPasswordTooShort
// -------------------------------
// Given a password that is too short (for example, fewer than 8 characters),
// the service should return ErrPasswordTooShort and avoid calling the repository.
func TestSignUp_WeakPasswordTooShort(t *testing.T) {
	svc := newTestAuthService()

	err := svc.SignUp(context.Background(), "test@example.com", "Abc123")

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrPasswordTooShort)
}

// TestSignUp_WeakPasswordNoDigit
// ------------------------------
// Given a password that has no digits (only letters), the service should
// return ErrPasswordTooWeak and avoid calling the repository.
func TestSignUp_WeakPasswordNoDigit(t *testing.T) {
	svc := newTestAuthService()

	err := svc.SignUp(context.Background(), "test@example.com", "OnlyLetters")

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrPasswordTooWeak)
}

// TestSignUp_OK
// -------------
// "Happy path" scenario:
//
// - Email and password satisfy the validation rules.
// - The service hashes the password.
// - The service calls repo.Create exactly once.
// - The stored password matches the hashed value, not the raw password.
func TestSignUp_OK(t *testing.T) {
	// Dedicated fakeUserRepo is used here so createdUsers can be inspected.
	repo := &fakeUserRepo{}
	svc := newTestAuthServiceWithFakeRepo(repo)

	email := "test@example.com"
	password := "Supersecret123" // valid: >= 8 chars, contains letters + digits

	// When: signing up with valid data
	err := svc.SignUp(context.Background(), email, password)

	// Then: SignUp should return no error
	require.NoError(t, err)

	// And: repo.Create must have been called exactly once
	require.Len(t, repo.createdUsers, 1)

	created := repo.createdUsers[0]

	// And: email is passed through unchanged
	assert.Equal(t, email, created.Email)

	// And: the user ID is non-empty (service generated one)
	assert.NotEmpty(t, created.ID)

	// And: the stored password is not the plain password,
	// but the hashed version produced by fakeHasher.
	assert.NotEqual(t, password, created.Password)
	assert.Equal(t, "HASHED:"+password, created.Password)
}
