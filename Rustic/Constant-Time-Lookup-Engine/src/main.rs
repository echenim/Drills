fn main() {
    let mut engine = FastLookupEngine::new();
    engine.put("apple".to_string());
    engine.put("banana".to_string());
    engine.put("cherry".to_string());

    println!("\n Initial state: {:?} \n", engine.data);

    engine.check("banana");
    engine.remove("banana");
    println!("\n current state: {:?} \n", engine.data);
    engine.check("banana");
    engine.check("apple");
    engine.check("cherry");
    engine.check("banana");
    println!("\n current state: {:?} \n", engine.data);
}

pub struct FastLookupEngine {
    data: Vec<String>,
    index: std::collections::HashMap<String, usize>,
}

impl FastLookupEngine {
    pub fn new() -> Self {
        FastLookupEngine {
            data: Vec::new(),
            index: std::collections::HashMap::new(),
        }
    }

    pub fn put(&mut self, value: String) {
        if !self.index.contains_key(&value) {
            self.data.push(value.clone());
            self.index.insert(value, self.data.len() - 1);
        }
    }

    pub fn check(&self, value: &str) {
        if self.index.contains_key(value) {
            println!("Yes, {} exists", value);
        } else {
            println!("No, {} does not exist", value);
        }
    }

    pub fn remove(&mut self, value: &str) {
        if let Some(&k) = self.index.get(value) {
            let last_index = self.data.len() - 1;
            self.data.swap(k, last_index);
            self.data.pop();
            self.index.remove(value);
            if k != last_index {
                let swapped_value = &self.data[k];
                self.index.insert(swapped_value.clone(), k);
            }
        }
    }
}
