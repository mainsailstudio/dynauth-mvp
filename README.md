# Please don't use this yet
### It's so not secure
And it probably won't work for you anyway! Check back in a week, maybe this notice will be gone and there will be install instructions instead... who knows? :)
# This is the Minimum Remarkable (Viable) Product for Dynauth
It contains a simple password manager that presents an API to a Chrome extension
## Packages
- pwm (PassWord Manager): Contains Golang code for the password manager API to perform simple CRUD operations on the password database as well as more complex security operations such as encryption, decryption, and "speading honey"
- dynauth (Dynamic Authentication): Contains Golang code for the dynamic authentication API to register and authenticate users
- chrome-extension: Contains mostly JavaScript code to provide a simple functional interface for using the password manager within a browser