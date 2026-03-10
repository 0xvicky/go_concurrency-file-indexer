/*
Error 1
An unbuffered channel requires both sides to be ready at the same time:

sender  ↔  receiver

Think of it like handing someone a package. You cannot drop it unless someone is there to receive it.

Right now the situation is:

scanner → trying to send
main → not receiving yet

Because main() is still waiting for scanner() to return.

So the send operation blocks.

The scanner waits forever for someone to receive.

But main() will only start receiving after scanner() finishes.

And scanner() cannot finish because it’s stuck sending.

You now have a perfect circular wait.

That is the deadlock.

=================================================

*/

/*
Program starts
      ↓
Create shared structures (channels, result map, mutex, waitgroup)
      ↓
Start worker goroutines
      ↓
Scanner walks through directory
      ↓
Scanner sends file paths into channel
      ↓
Workers continuously read file paths
      ↓
Workers open file and compute hash
      ↓
Workers safely store result
      ↓
Scanner finishes → channel closed
      ↓
Workers detect closed channel → exit
      ↓
WaitGroup waits for all workers
      ↓
Program prints results
      ↓
Program exits
*/

/*
Filesystem✅
    ↓
Directory Scanner ✅
    ↓
File Path Channel ✅
    ↓
Worker Pool (goroutines) ✅
    ↓
Hash Computation
    ↓
Results Map (mutex protected)
    ↓
Program Summary
*/
