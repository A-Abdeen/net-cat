## NET-CAT

This project is a TCP Chat Server that allows users to participate in a group chat. To run the server and connect as a client, follow these steps:

**Running the Server**

1. Clone or download this repository to your local machine.
2. Navigate to the project directory and enter the server directory.
3. Open a terminal or command prompt.

### To start the server, run the following command:

```bash
go run .
```

The server will start listening on the default port :8989. You can specify a different port if needed, like this:

```bash
go run . 2525
```

### Connecting as a client

To connect to the server as a client, use the nc (NetCat) command followed by the server's IP address and port number. For example:

```bash
nc <SERVER_IP> <PORT>
```

- Replace <SERVER_IP> with the actual IP address of the server and <PORT> with the port number it's listening on. Once connected, you can start chatting with other clients in the group chat.

- When prompted provide a username and choose the group chat you want to join.

That's it! You're now ready to run the TCP Chat Server and connect as a client to join the chat.

-Note: flags available: '--name': allows you to change your name in the group chat.(ex: --name Ahmed)
                        '--users': shows you the number of users connected.
                        '--switch': allows you to switch from one group chat to another.(ex: --switch adnan)

**Authors**

- Adnan Jaberi
- Ahmad alali
- ahmad abdeen
