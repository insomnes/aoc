use std::str::FromStr;

pub fn load_input(filename: &str) -> Vec<String> {
    std::fs::read_to_string(filename)
        .expect("Failed to read file")
        .lines()
        .map(|line| line.to_string())
        .collect()
}

pub enum Direction {
    Forward(i32),
    Down(i32),
    Up(i32),
}

impl FromStr for Direction {
    type Err = String;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let split = s.split(' ').collect::<Vec<&str>>();
        if split.len() != 2 {
            return Err("Wrong number of arguments".to_string());
        }
        let (dir, num) = (split[0], split[1]);
        let num = num.parse().expect("Failed to parse number");
        match dir {
            "forward" => Ok(Direction::Forward(num)),
            "down" => Ok(Direction::Down(num)),
            "up" => Ok(Direction::Up(num)),
            _ => Err("Wrong direction".to_string()),
        }
    }
}

pub fn parse_input_part_one(input: &[String]) -> Vec<Direction> {
    input.iter().map(|line| line.parse().unwrap()).collect()
}

pub fn parse_input_part_two(input: &[String]) -> Vec<Direction> {
    parse_input_part_one(input)
}

pub fn part_one(input: Vec<Direction>) -> i32 {
    let mut x = 0;
    let mut y = 0;
    for cmd in input.iter() {
        match cmd {
            Direction::Forward(num) => x += num,
            Direction::Down(num) => y += num,
            Direction::Up(num) => {
                y -= num;
            }
        }
    }
    x * y
}

pub fn part_two(input: Vec<Direction>) -> i32 {
    let mut x = 0;
    let mut y = 0;
    let mut aim = 0;
    for cmd in input.iter() {
        match cmd {
            Direction::Forward(num) => {
                x += num;
                y += num * aim;
            }
            Direction::Down(num) => aim += num,
            Direction::Up(num) => aim -= num,
        }
    }
    x * y
}
