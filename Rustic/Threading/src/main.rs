use std::thread::spawn;
use std::thread::JoinHandle;

pub fn test_threads() {
    let mut x = 0u128;

    for i in 0..100_000_000_000 {
        x += i;
    }
    println!("x: {}", x);
}

pub fn spawn_thread() {
    let thread_fn = || {
        test_threads();
    };
    let handle: JoinHandle<()> = spawn(thread_fn);
    handle.join().unwrap();
}

fn main() {
    spawn_thread();
}
