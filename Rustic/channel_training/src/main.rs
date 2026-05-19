use std::sync::mpsc;
use std::thread;

// fn send_message(tx: mpsc::Sender<&str>) {
//     tx.send("Hello from the channel!").unwrap();
// }

fn send_multi_message(tx: mpsc::Sender<&str>) {
    let messages = vec!["Hello", "from", "the", "channel!"];
    for msg in messages {
        tx.send(msg).unwrap();
        thread::sleep(std::time::Duration::from_secs(1));
    }
}

fn send_multi_message_x(tx: mpsc::Sender<&str>) {
    let messages = vec!["More", "messages", "for", "You"];
    for msg in messages {
        tx.send(msg).unwrap();
        thread::sleep(std::time::Duration::from_secs(2));
    }
}

// fn receive_message(rx: mpsc::Receiver<&str>) {
//     let received = rx.recv().unwrap();
//     println!("Received: {}", received);
// }

fn receive_multi_message(rx: mpsc::Receiver<&str>) {
    for received in rx {
        println!("Received: {}", received);
    }
}

fn main() {
    let (tx, rx) = mpsc::channel();
    let tx1 = tx.clone();
    // thread::spawn(|| {
    //     send_message(tx);
    // });
    // receive_message(rx);

    thread::spawn(|| {
        send_multi_message(tx);
    });
    thread::spawn(|| {
        send_multi_message_x(tx1);
    });
    receive_multi_message(rx);
}
