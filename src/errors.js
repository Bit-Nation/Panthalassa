// @flow

class InvalidPrivateKeyError extends Error{}

class PasswordMismatch extends Error{}

/**
 * Is thrown if a method on the secure storage implementation is missing
 */
class UnsatisfiedSecureStorageImplementationError extends Error{

    constructor(missingMethodName: string) {

        super('Missing method: "'+missingMethodName+'" in secure storage implementation');

    }

}

module.exports = {
    InvalidPrivateKeyError,
    UnsatisfiedSecureStorageImplementationError,
    PasswordMismatch
};
