## Declaration
This system is an example of using tcpx.It advises but not force users to design your system like this.

## Server brokers
Server are divided into brokers: center, register, userPool. They works below.

#### Center
Center handles all events from clients.It can not only scale out horizontally based on a certain event, but also scale out for varies events.

To interfere with user.It will grasp userInfo from pool and thus grasp which pool this user is in. Then It will build a connection(called bridge) to the specifc pool.

- Receiving message from client.
- Connect to specific userPool for further operaion.

#### Register
Register works for registering user info and pool info.It also works as a redis proxy.All interfering actions with redis can be handled here.

- Storage and provide user login info.
- Storage and provide pool info.

#### UserPool
UserPool can scale out without number limit. Once a user client send online to a user-pool broker, then this user will be joint with it.

- Save user connections

**Most Important Point:**
All server brokers can easily scale out without side effect.

## Broker doc
MessageID for different brokers are designed first. For broker userPool, messageID ranges 1-100. For center broker, messageID ranges 101-200.


