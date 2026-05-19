use crossbeam_channel::{bounded, Receiver, Sender};
use std::sync::{
    atomic::{AtomicBool, Ordering},
    Arc,
};
use std::thread;

// 1. log() must be thread-safe.
// 2. At most X messages can be processed concurrently.
// 3. At most Y messages can wait in the queue.
// 4. If the queue is full, log() must block.
// 5. After shutdown(), new log() calls should return immediately.
// 6. shutdown() must wait until all queued messages are processed.

pub struct Logger {
    sender: Sender<String>,
    shutdown_flag: Arc<AtomicBool>,
    worker_handle: Vec<thread::JoinHandle<()>>,
}

impl Logger {
    pub fn new(max_concurrent: usize, max_queue: usize) -> Self {
        let (sender, receiver) = bounded(max_queue);
        let shutdown_flag = Arc::new(AtomicBool::new(false));
        let mut worker_handle = Vec::with_capacity(max_concurrent);

        for _ in 0..max_concurrent {
            let receiver = receiver.clone();
            let shutdown_flag = Arc::clone(&shutdown_flag);
            let handle = thread::spawn(move || {
                loop {
                    crossbeam::select! {
                        recv(receiver) -> msg => {
                            match msg {
                                Ok(message) => {
                                    println!("Logged: {}", message);
                                },
                                Err(_) => break, // Channel closed
                            }
                        }
                    }
                    if shutdown_flag.load(Ordering::SeqCst) {
                        while let Ok(message) = receiver.try_recv() {
                            println!("Logged: {}", message);
                        }
                        break;
                    }
                }
            });
            worker_handle.push(handle);
        }

        drop(receiver);

        Logger {
            sender,
            shutdown_flag,
            worker_handle,
        }
    }

    pub fn log(&self, message: String) {
        if self.shutdown_flag.load(Ordering::SeqCst) {
            return;
        }
        self.sender.send(message).unwrap();
    }

    pub fn shutdown(self) {
        self.shutdown_flag.store(true, Ordering::SeqCst);
        drop(self.sender.clone());
        for handle in self.worker_handle {
            let _ = handle.join();
        }
    }
}

fn main() {
    let logger = Logger::new(3, 5);

    for i in 0..100_000_000 {
        let msg = format!("Message {}", i);
        logger.log(msg);
    }

    logger.shutdown();
}
