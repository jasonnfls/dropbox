# dropbox
A dead-simple file exchanger written in Golang built upon RESTful HTTP. Built with simplicity, everything is within one single source file for readability and maintainability.

Have you ever found yourself in the following situations constantly?
* Signing up a dodgy website like yousendit.com just to send a 5MB PDF to a colleage you barely know.
* "lol can you email me that file? my email is ...." to a friend who keeps forgetting your email
* "I am sorry but I could not open the Google Drive link you shared previosly. The error was no permission. Would you please check ..." to a supervisor or someone you aren't first-name basis to

Well, you need this program.

## features
1. Directory listing/navigation
2. Downloading/Deletion of files
3. Creation/Deletion of (sub)directories

## configuration snippet for nginx reverse-proxy
'''
server {
    client_max_body_size 2048M;
    charset   utf-8;
    location / {
        proxy_pass http://127.0.0.1:1234;
    }
}
'''
