# Please don't use this yet
### It's so not secure
And it probably won't work for you anyway! Check back in a week, maybe this notice will be gone and there will be install instructions instead... who knows? :)

# Dynauth: Password Manager (MVP)

#### Created by Connor Peters, looking for collaboration. Please contact me via: https://dynauth.io/contact

## Introduction
The overarching goal is to provide seamless, secure, and self-owned identity to everyone. To keep this short, the first step in accomplishing that is to create a *more* secure authentication process to users on the Internet by making a password manager utilizing dynamic authentication. Read more about the password manager here: https://dynauth.io/password-manager and read more about dynamic authentication here: https://dynauth.io/dynamic-authentication

## Design
There are of course, a billion amazing password managers out there that serve the needs of many quite well. But all of them use a master passwords as the final barrier to entry - Dynauth is different.

Let's start with some important design decisions that differentiates this password manager from others.

1. It's a self-owned locally hosted database with an API
   1. The most important aspect of Dynauth is that everything is self-owned. That means having a hosted cloud service (like LastPass) is off the table, despite the great user experience it provides. So instead, users will use Dynauth to simply sync data between their devices seamlessly, with Dynauth never storing any information during the process. That means that the core password manager has to be self-hosted locally, but present an API that Dynauth can take advantage of to improve the user experience.
   2. This also enables Dynauth to build front-end apps to improve the user experience, things such as browser extensions that bypass the login form entirely and post the data directly to the service to authenticate the user. See the chrome extension folder for implementation details there.
2. The passwords need to be securely encrypted. Obviously.
   1. That means we need a method to encrypt the passwords, but then decrypt them upon user authentication. Typically from my research, this is done with PBE (Password Based Encryption), a form of symmetric-key encryption that pairs wonderfully with a master password as the symmetric key. One problem with this though... Dynauth does not have a master password, by design. That means the normal method of using PBKDF2 with SHA256 won't work. We'll need a better alternative (under development).
   2. Each password will be encrypted with something call honey encryption. Basically, that means that an attacker performing an offline attack will never be rejected at authentication, rather, the decrypted passwords that come as a result are just fake decoys provided to force the attacker to slow down to online attack speeds. More info here: https://en.wikipedia.org/wiki/Honey_encryption

## API Functions
Since the core password manager is locally hosted, providing an API, we need to determine what exactly the API presents and how it does so securely.

Let's just start by throwing ideas at the wall and listing all the functions necessary to make a functioning password manager. This is just a basic outline for ideation, everything will be encrypted/decrypted properly.

- createPassword(): POST to insert new password
- getPassword(): GET to get password based on URL
- updatePassword(): PUT that updates existing password
- deletePassword():	DELETE that deletes an existing password

## Conclusion
This is of course a work in progress, please don't use any code here seriously yet. If you're interested in this, please contact me via https://dynauth.io/contact.

Contact me if you:
1. Want to work with me
2. Want to give me some helpful input
3. Want to criticize me
4. Want to yell at me and make me question everything I stand for
5. Want a risky investment opportunity