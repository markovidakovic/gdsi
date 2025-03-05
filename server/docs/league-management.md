# Tennis League Management System

Lets break down how this system would work from start to finish

## ELO Rating Fundamentals

- All players start with a base ELO rating (typically 1200)
- New players are marked as "provisional" until they complete 10 matches
- Provisional players ratings change more dramatically (using a higher K-factor) 
- Elo adjusts after each match based on the outcome and opponent's rating

## League Structure

- Leagues are organized in a pyramid structure
- Multiple tiers (1,2,3, etc.) with groups within each tier (a,b,c)
- Higher tier = more skilled players (tier 1 is the highest)
- Each league constains 4-6 players


## Monthly Season Cycle

### End of Current Season:

1. Final matches are completed
2. League standings are calculated based on match results
3. For each player, the system records:
   - Their final league (previous_league_id)
   - Their final position (previous_rank)
4. The system determines who gets promoted/relegated:
   - Top 2 players move up one tier
   - Bottom 2 players move down one tier
   - Middle players remain in the same tier

### Creating a New Season:

1. Admin initiates new season creation
2. System creates a fresh league structure
3. Players are assigned to leagues based on:
   - Previous performance (promotion/relegation)
   - Current ELO rating (especially for new players)
4. The system handles edge cases:
   - Already in top tier (can't be promoted further)
   - Already in bottom tier (can't be relegated further)
   - Leagues that are too full or too empty

## Match Processing

1. Players schedule and play matches
2. After each match:
   - System records the result
   - Updates both players ELO ratings
   - Updates league standings
   - Checks if provisional status should be removed

## Ongoing League Management

- Players can see their current league, position, and ELO
- Admin can monitor league health (balance, activity)
- System maintains league integrity through the promotion/relegation mechanism

This creates a self-balancing competitive system where:

- Players compete against others of similar skill
- Improvement is rewarded with promotion
- ELO provides an objective measure of skill
- The montly reset keeps competition fresh