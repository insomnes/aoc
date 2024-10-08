type ParsedInput = [usize; 2];
use std::collections::HashMap;

pub fn load_input(filename: &str) -> Vec<String> {
    std::fs::read_to_string(filename)
        .expect("Failed to read file")
        .lines()
        .map(|line| line.to_string())
        .collect()
}

pub fn parse_input_part_one(input: &[String]) -> ParsedInput {
    assert_eq!(input.len(), 2);
    let p1 = input[0].split(": ").collect::<Vec<&str>>()[1]
        .parse()
        .unwrap();
    let p2 = input[1].split(": ").collect::<Vec<&str>>()[1]
        .parse()
        .unwrap();
    [p1, p2]
}

pub fn parse_input_part_two(input: &[String]) -> ParsedInput {
    parse_input_part_one(input)
}

pub fn part_one(input: ParsedInput) -> usize {
    let (p1, p2) = find_repeating_postions(input);
    let mut p1_rounds = find_player_win_round(&p1);
    let mut p2_rounds = find_player_win_round(&p2);
    println!("P1: {} P2: {}", p1_rounds, p2_rounds);
    if p1_rounds < p2_rounds {
        // 3 rolls per round, 2 players, second player didn't get a final role so -3
        let total_rolls = p1_rounds * 3 * 2 - 3;
        p2_rounds = p1_rounds - 1;
        let mut p2_score = p2.iter().sum::<usize>() * (p2_rounds / p2.len());
        if p2_rounds % p2.len() != 0 {
            p2_score += p2[..p2_rounds % p2.len()].iter().sum::<usize>();
        }
        println!("P2 score: {}, Total rolls: {}", p2_score, total_rolls);

        return total_rolls * p2_score;
    }

    let total_rolls = p2_rounds * 3 * 2;
    p1_rounds = p2_rounds;
    let mut p1_score = p1.iter().sum::<usize>() * (p1_rounds / p1.len());
    if p1_rounds % p1.len() != 0 {
        p1_score += p1[..p1_rounds % p1.len()].iter().sum::<usize>();
    }
    println!("P1 score: {}, Total rolls: {}", p1_score, total_rolls);
    total_rolls * p1_score
}

fn find_player_win_round(all_pos: &[usize]) -> usize {
    let p_sum: usize = all_pos.iter().sum();
    let p_fulls = 1000 / p_sum * all_pos.len();
    let p_remainder = 1000 % p_sum;
    if p_remainder == 0 {
        return p_fulls;
    }

    let mut sum = 0;
    let mut rounds = None;
    for (i, score) in all_pos.iter().enumerate() {
        sum += score;
        if sum >= p_remainder {
            rounds = Some(i + 1);
            break;
        }
    }
    if let Some(rounds) = rounds {
        p_fulls + rounds
    } else {
        panic!("No solution found");
    }
}

fn find_repeating_postions(input: ParsedInput) -> (Vec<usize>, Vec<usize>) {
    let mut p1_pos = input[0];
    let mut p1_all = Vec::new();
    let mut p1_found = 0;

    let mut p2_pos = input[1];
    let mut p2_all = Vec::new();
    let mut p2_found = 0;

    for n in 1..100 {
        if p1_found > 0 && p2_found > 0 {
            break;
        };
        if n % 2 != 0 {
            if p1_found > 0 {
                continue;
            }
            p1_pos = move_player(p1_pos, n);
            p1_all.push(p1_pos);
            if check_player_position(&p1_all) {
                p1_found = n / 2;
                p1_all = p1_all[..p1_all.len() / 2].to_vec();
                println!("{} ({}) P1: {:?}", n, p1_found, p1_all);
            }
            continue;
        };

        if p2_found > 0 {
            continue;
        }
        p2_pos = move_player(p2_pos, n);
        p2_all.push(p2_pos);
        if check_player_position(&p2_all) {
            p2_found = n / 2;
            p2_all = p2_all[..p2_all.len() / 2].to_vec();
            println!("{} ({}) P2: {:?}", n, p2_found, p2_all);
        }
    }

    (p1_all, p2_all)
}

fn check_player_position(all: &[usize]) -> bool {
    if all.len() < 2 || all.len() % 2 != 0 {
        return false;
    }
    all[..all.len() / 2] == all[all.len() / 2..]
}

fn move_player(pos: usize, roll: usize) -> usize {
    let new_pos = pos + calculate_roll_sum(roll);
    if new_pos % 10 == 0 {
        return 10;
    }
    new_pos % 10
}

fn calculate_roll_sum(r_num: usize) -> usize {
    let a = fix_100_overflow(3 * (r_num - 1)) + 1;
    let b = fix_100_overflow(3 * (r_num - 1) + 1) + 1;
    let c = fix_100_overflow(3 * (r_num - 1) + 2) + 1;

    a + b + c
}

fn fix_100_overflow(num: usize) -> usize {
    if num < 100 {
        return num;
    }
    num % 100
}

const POSSIBLE_ROLLS: [usize; 7] = [3, 4, 5, 6, 7, 8, 9];
const ROLL_FREQ: [usize; 10] = [0, 0, 0, 1, 3, 6, 7, 6, 3, 1];
const WIN_POS_2: usize = 21;

pub fn part_two(input: ParsedInput) -> usize {
    let (p1, p2) = (input[0], input[1]);

    let mut cache = HashMap::new();
    let (p1_wins, p2_wins) = count_dirac_wins(p1, p2, 0, 0, &mut cache);
    let max_wins = p1_wins.max(p2_wins);
    println!("P1: {}, P2: {}, Max: {}", p1_wins, p2_wins, max_wins);
    max_wins
}

fn count_dirac_wins(
    cur: usize,
    other: usize,
    cur_score: usize,
    other_score: usize,
    cache: &mut HashMap<(usize, usize, usize, usize), (usize, usize)>,
) -> (usize, usize) {
    // We can save a lot by memoization since we have a lot of repeated states
    if let Some(&wins) = cache.get(&(cur, other, cur_score, other_score)) {
        return wins;
    }
    if cur_score >= WIN_POS_2 {
        return (1, 0);
    }
    if other_score >= WIN_POS_2 {
        return (0, 1);
    }
    let mut cur_wins = 0;
    let mut other_wins = 0;

    for roll in POSSIBLE_ROLLS {
        let new_pos = (cur + roll - 1) % 10 + 1;
        let new_score = cur_score + new_pos;
        // The tricky part here is to switch the players
        let (new_other_wins, new_cur_wins) =
            count_dirac_wins(other, new_pos, other_score, new_score, cache);
        cur_wins += new_cur_wins * ROLL_FREQ[roll];
        other_wins += new_other_wins * ROLL_FREQ[roll];
    }
    cache.insert((cur, other, cur_score, other_score), (cur_wins, other_wins));
    (cur_wins, other_wins)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_fix_100_overflow() {
        assert_eq!(fix_100_overflow(2), 2);
        assert_eq!(fix_100_overflow(100), 0);
        assert_eq!(fix_100_overflow(101), 1);
        assert_eq!(fix_100_overflow(202), 2);
    }

    #[test]
    fn test_calculate_roll_sum() {
        assert_eq!(calculate_roll_sum(1), 6);
        assert_eq!(calculate_roll_sum(2), 15);
        assert_eq!(calculate_roll_sum(8), 69);
        assert_eq!(calculate_roll_sum(33), 294);
        assert_eq!(calculate_roll_sum(34), 103);
        assert_eq!(calculate_roll_sum(35), 12);
        // 99 + 100 + 1 = 200
        assert_eq!(calculate_roll_sum(67), 200);
        assert_eq!(calculate_roll_sum(68), 9);
        // 95 + 96 + 97 = 288
        assert_eq!(calculate_roll_sum(99), 288);
        // 98 + 99 + 100 = 297
        assert_eq!(calculate_roll_sum(100), 297);
        assert_eq!(calculate_roll_sum(101), 6);
    }

    #[test]
    fn test_check_player_position() {
        assert!(!check_player_position(&vec![]));
        assert!(check_player_position(&vec![1, 2, 3, 1, 2, 3]));
        assert!(!check_player_position(&vec![1, 2, 3, 1, 2]));
        assert!(check_player_position(&vec![5, 2, 1, 1, 5, 2, 1, 1]));
    }
}
