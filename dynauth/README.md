# Dyanuth: Dynamic Authentication

#### Created by Connor Peters, looking for collaboration. Please contact me via: https://dynauth.io/contact

## Introduction
Dynamic authentication is designed to completely replace passwords on the Internet. It's not supposed to be a perfect authentication scheme, but it does solve lot's of the problems associated with passwords. Read more about the technical details behind dynamic authentication here: https://dynauth.io/dynamic-authentication

## Design
Since dynamic authentication is designed to replace passwords and aid the effort of providing online identity, it must be designed similarly to passwords to be able to function in the same environment.

## API Functions
- getLocks(): GET that retrieves a user's locks and temporarily stores requested locks on the server side
  - GET: /auth/users/locks?userid=123456789
- authenticate(): GET that compares the submitted hash to storage to determine if the user is authenticated. Returns a token
  - GET: /auth/users/keys?userid=123456789&secret=hashedsecrethere
- register(): POST that registers a user using email and dynamic authentication
  - POST: /auth/users?email=example@example.com&secret=allhashedlocksandkeyshere
- logAuthenticationAttempt()
  - POST /auth/users/logs?userid=123456789&logVar=example

## Conclusion
This is of course a work in progress, please don't use any code here seriously yet. If you're interested in this, please contact me via https://dynauth.io/contact.

Contact me if you:
1. Want to work with me
2. Want to give me some helpful input
3. Want to criticize me
4. Want to yell at me and make me question everything I stand for
5. Want a risky investment opportunity