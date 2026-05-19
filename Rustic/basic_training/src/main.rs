use std::thread;

pub struct Shape {
    width: f32,
    height: f32,
}

impl Shape {
    fn area_of_square(&self) -> f32 {
        self.width * self.height
    }

    fn area_of_rectangle(&self) -> f32 {
        self.width * self.height
    }
    fn area_of_triangle(&self) -> f32 {
        (self.width * self.height) as f32 / 2.0
    }
    fn area_of_trapezoid(&self, base1: f32, base2: f32) -> f32 {
        ((base1 + base2) as f32 / 2.0) * self.height as f32
    }
}

fn fetch_user(area: Vec<f32>) {
    println!("Area: {:?}", area);
}

fn main() {
    println!("Hello, world!");
    let mut area = Vec::new();
    area.push(
        Shape {
            width: 5.0,
            height: 5.0,
        }
        .area_of_square(),
    );
    area.push(
        Shape {
            width: 5.0,
            height: 10.0,
        }
        .area_of_rectangle(),
    );
    area.push(
        Shape {
            width: 5.0,
            height: 10.0,
        }
        .area_of_triangle(),
    );
    area.push(
        Shape {
            width: 5.0,
            height: 10.0,
        }
        .area_of_trapezoid(5.0, 10.0),
    );
    let handle = thread::spawn(|| {
        fetch_user(area);
    });

    handle.join().unwrap();
}
