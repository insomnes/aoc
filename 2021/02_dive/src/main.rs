const TEST_PATH: &str = "test.txt";
const INPUT_PATH: &str = "input.txt";

use dive::{load_input, parse_input_part_one, parse_input_part_two, part_one, part_two};

fn main() {
    let mut fp = TEST_PATH;
    let args = std::env::args().collect::<Vec<String>>();
    if args.len() > 1 {
        fp = INPUT_PATH;
    }
    let inp = load_input(fp);

    {
        let mut start = std::time::Instant::now();

        let parsed = parse_input_part_one(&inp);
        let parsed_elapsed = start.elapsed();

        start = std::time::Instant::now();

        let result = part_one(parsed);
        let elapsed = start.elapsed();

        println!("Result one: {:?}", result);
        println!("Elapsed one: {:?}, parsing: {:?}", elapsed, parsed_elapsed);
    }
    {
        let mut start = std::time::Instant::now();

        let parsed = parse_input_part_two(&inp);
        let parsed_elapsed = start.elapsed();

        start = std::time::Instant::now();

        let result = part_two(parsed);

        let elapsed = start.elapsed();

        println!("Result two: {:?}", result);
        println!("Elapsed two: {:?}, parsing: {:?}", elapsed, parsed_elapsed);
    }
}
