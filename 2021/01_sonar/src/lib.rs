pub fn load_input(filename: &str) -> Vec<String> {
    std::fs::read_to_string(filename)
        .expect("Failed to read file")
        .lines()
        .map(|line| line.to_string())
        .collect()
}

pub fn parse_input_part_one(input: &[String]) -> Vec<i32> {
    input.iter().map(|line| line.parse().unwrap()).collect()
}

pub fn parse_input_part_two(input: &[String]) -> Vec<i32> {
    parse_input_part_one(input)
}

pub fn part_one(input: Vec<i32>) -> i32 {
    let mut increased = 0;
    input.windows(2).for_each(|pair| {
        if pair[1] > pair[0] {
            increased += 1;
        }
    });

    increased
}

pub fn part_two(input: Vec<i32>) -> i32 {
    let mut increased = 0;
    let mut triplet_iter = input.windows(3);
    let mut prev = triplet_iter.next().expect("Empty iter");
    for next in triplet_iter {
        if next.iter().sum::<i32>() - prev.iter().sum::<i32>() > 0 {
            increased += 1;
        }
        prev = next;
    }
    increased
}
