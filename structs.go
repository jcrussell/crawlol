package main

type Summoner struct {
	Id            int64  // Summoner ID.
	Name          string // Summoner name.
	ProfileIconId int    // ID of the summoner icon associated with the summoner.
	RevisionDate  int64  // Date summoner was last modified specified as epoch milliseconds.
	SummonerLevel int64  // Summoner level associated with the summoner.

	LastCrawled int64 // Last time summoner's games were crawled, specified as epoch milliseconds.
}

// There are a lot more fields returned by the API request, however, we really
// only care about the MatchID since we'll use it to request the full details.
type MatchSummary struct {
	MatchID int64 // ID of the match
}

type MatchDetail struct {
	MapID int // Match map ID
	// Match creation time. Designates when the team select lobby is created
	// and/or the match is made through match making, not when the game actually
	// starts.
	MatchCreation int64
	MatchDuration int64 // Match duration
	MatchID       int64 // ID of the match
	// Match mode (legal values: CLASSIC, ODIN, ARAM, TUTORIAL, ONEFORALL,
	// ASCENSION, FIRSTBLOOD)
	MatchMode string
	// Match type (legal values: CUSTOM_GAME, MATCHED_GAME, TUTORIAL_GAME)
	MatchType             string
	MatchVersion          string                // Match version
	ParticipantIdentities []ParticipantIdentity // Participant identity information
	Participants          []Participant         // Participant information
	// Match queue type (legal values: CUSTOM, NORMAL_5x5_BLIND, RANKED_SOLO_5x5,
	// RANKED_PREMADE_5x5, BOT_5x5, NORMAL_3x3, RANKED_PREMADE_3x3,
	// NORMAL_5x5_DRAFT, ODIN_5x5_BLIND, ODIN_5x5_DRAFT, BOT_ODIN_5x5,
	// BOT_5x5_INTRO, BOT_5x5_BEGINNER, BOT_5x5_INTERMEDIATE, RANKED_TEAM_3x3,
	// RANKED_TEAM_5x5, BOT_TT_3x3, GROUP_FINDER_5x5, ARAM_5x5, ONEFORALL_5x5,
	// FIRSTBLOOD_1x1, FIRSTBLOOD_2x2, SR_6x6, URF_5x5, BOT_URF_5x5,
	// NIGHTMARE_BOT_5x5_RANK1, NIGHTMARE_BOT_5x5_RANK2, NIGHTMARE_BOT_5x5_RANK5,
	// ASCENSION_5x5)
	QueueType string
	Region    string // Region where the match was played
	// Season match was played (legal values: PRESEASON3, SEASON3, PRESEASON2014,
	// SEASON2014)
	Season   string
	Teams    []Team   // Team information
	Timeline Timeline // Match timeline data.
}

type Participant struct {
	ChampionID            int                 // Champion ID
	Masteries             []Mastery           // List of mastery information
	ParticipantIdentities int                 // Participant ID
	Runes                 []Rune              // List of rune information
	Spell1ID              int                 // First summoner spell ID
	Spell2ID              int                 // Second summoner spell ID
	Stats                 ParticipantStats    // Participant statistics
	TeamID                int                 // Team ID
	Timeline              ParticipantTimeline // Timeline data
}

type ParticipantIdentity struct {
	ParticipantID int    // Participant ID
	Player        Player // Player information
}

type Player struct {
	MatchHistoryUri string // Match history URI
	ProfileIcon     int    // Profile icon ID
	SummonerID      int64  // Summoner ID
	SummonerName    string // Summoner name
}

type Team struct {
	Bans                 []BannedChampion // If game was draft mode, contains banned champion data, otherwise null
	BaronKills           int              // Number of times the team killed baron
	DominionVictoryScore int64            // If game was a dominion game, specifies the points the team had at game end, otherwise null
	DragonKills          int              // Number of times the team killed dragon
	FirstBaron           bool             // Flag indicating whether or not the team got the first baron kill
	FirstBlood           bool             // Flag indicating whether or not the team got first blood
	FirstDragon          bool             // Flag indicating whether or not the team got the first dragon kill
	FirstInhibitor       bool             // Flag indicating whether or not the team destroyed the first inhibitor
	FirstTower           bool             // Flag indicating whether or not the team destroyed the first tower
	InhibitorKills       int              // Number of inhibitors the team destroyed
	TeamID               int              // Team ID
	TowerKills           int              // Number of towers the team destroyed
	VilemawKills         int              // Number of times the team killed vilemaw
	Winner               bool             // Flag indicating whether or not the team won
}

type Timeline struct {
	FrameInterval int64   // Time between each returned frame in milliseconds.
	Frames        []Frame // List of timeline frames for the game.
}

type ParticipantStats struct {
	Assists                         int64 // Number of assists
	ChampLevel                      int64 // Champion level achieved
	CombatPlayerScore               int64 // If game was a dominion game, player's combat score, otherwise 0
	Deaths                          int64 // Number of deaths
	DoubleKills                     int64 // Number of double kills
	FirstBloodAssist                bool  // Flag indicating if participant got an assist on first blood
	FirstBloodKill                  bool  // Flag indicating if participant got first blood
	FirstInhibitorAssist            bool  // Flag indicating if participant got an assist on the first inhibitor
	FirstInhibitorKill              bool  // Flag indicating if participant destroyed the first inhibitor
	FirstTowerAssist                bool  // Flag indicating if participant got an assist on the first tower
	FirstTowerKill                  bool  // Flag indicating if participant destroyed the first tower
	GoldEarned                      int64 // Gold earned
	GoldSpent                       int64 // Gold spent
	InhibitorKills                  int64 // Number of inhibitor kills
	Item0                           int64 // First item ID
	Item1                           int64 // Second item ID
	Item2                           int64 // Third item ID
	Item3                           int64 // Fourth item ID
	Item4                           int64 // Fifth item ID
	Item5                           int64 // Sixth item ID
	Item6                           int64 // Seventh item ID
	KillingSprees                   int64 // Number of killing sprees
	Kills                           int64 // Number of kills
	LargestCriticalStrike           int64 // Largest critical strike
	LargestKillingSpree             int64 // Largest killing spree
	LargestMultiKill                int64 // Largest multi kill
	MagicDamageDealt                int64 // Magical damage dealt
	MagicDamageDealtToChampions     int64 // Magical damage dealt to champions
	MagicDamageTaken                int64 // Magic damage taken
	MinionsKilled                   int64 // Minions killed
	NeutralMinionsKilled            int64 // Neutral minions killed
	NeutralMinionsKilledEnemyJungle int64 // Neutral jungle minions killed in the enemy team's jungle
	NeutralMinionsKilledTeamJungle  int64 // Neutral jungle minions killed in your team's jungle
	NodeCapture                     int64 // If game was a dominion game, number of node captures
	NodeCaptureAssist               int64 // If game was a dominion game, number of node capture assists
	NodeNeutralize                  int64 // If game was a dominion game, number of node neutralizations
	NodeNeutralizeAssist            int64 // If game was a dominion game, number of node neutralization assists
	ObjectivePlayerScore            int64 // If game was a dominion game, player's objectives score, otherwise 0
	PentaKills                      int64 // Number of penta kills
	PhysicalDamageDealt             int64 // Physical damage dealt
	PhysicalDamageDealtToChampions  int64 // Physical damage dealt to champions
	PhysicalDamageTaken             int64 // Physical damage taken
	QuadraKills                     int64 // Number of quadra kills
	SightWardsBoughtInGame          int64 // Sight wards purchased
	TeamObjective                   int64 // If game was a dominion game, number of completed team objectives (i.e., quests)
	TotalDamageDealt                int64 // Total damage dealt
	TotalDamageDealtToChampions     int64 // Total damage dealt to champions
	TotalDamageTaken                int64 // Total damage taken
	TotalHeal                       int64 // Total heal amount
	TotalPlayerScore                int64 // If game was a dominion game, player's total score, otherwise 0
	TotalScoreRank                  int64 // If game was a dominion game, team rank of the player's total score (e.g., 1-5)
	TotalTimeCrowdControlDealt      int64 // Total dealt crowd control time
	TotalUnitsHealed                int64 // Total units healed
	TowerKills                      int64 // Number of tower kills
	TripleKills                     int64 // Number of triple kills
	TrueDamageDealt                 int64 // True damage dealt
	TrueDamageDealtToChampions      int64 // True damage dealt to champions
	TrueDamageTaken                 int64 // True damage taken
	UnrealKills                     int64 // Number of unreal kills
	VisionWardsBoughtInGame         int64 // Vision wards purchased
	WardsKilled                     int64 // Number of wards killed
	WardsPlaced                     int64 // Number of wards placed
	Winner                          bool  // Flag indicating whether or not the participant won
}

type ParticipantTimeline struct {
	AncientGolemAssistsPerMinCounts ParticipantTimelineData // Ancient golem assists per minute timeline counts
	AncientGolemKillsPerMinCounts   ParticipantTimelineData // Ancient golem kills per minute timeline counts
	AssistedLaneDeathsPerMinDeltas  ParticipantTimelineData // Assisted lane deaths per minute timeline data
	AssistedLaneKillsPerMinDeltas   ParticipantTimelineData // Assisted lane kills per minute timeline data
	BaronAssistsPerMinCounts        ParticipantTimelineData // Baron assists per minute timeline counts
	BaronKillsPerMinCounts          ParticipantTimelineData // Baron kills per minute timeline counts
	CreepsPerMinDeltas              ParticipantTimelineData // Creeps per minute timeline data
	CsDiffPerMinDeltas              ParticipantTimelineData // Creep score difference per minute timeline data
	DamageTakenDiffPerMinDeltas     ParticipantTimelineData // Damage taken difference per minute timeline data
	DamageTakenPerMinDeltas         ParticipantTimelineData // Damage taken per minute timeline data
	DragonAssistsPerMinCounts       ParticipantTimelineData // Dragon assists per minute timeline counts
	DragonKillsPerMinCounts         ParticipantTimelineData // Dragon kills per minute timeline counts
	ElderLizardAssistsPerMinCounts  ParticipantTimelineData // Elder lizard assists per minute timeline counts
	ElderLizardKillsPerMinCounts    ParticipantTimelineData // Elder lizard kills per minute timeline counts
	GoldPerMinDeltas                ParticipantTimelineData // Gold per minute timeline data
	InhibitorAssistsPerMinCounts    ParticipantTimelineData // Inhibitor assists per minute timeline counts
	InhibitorKillsPerMinCounts      ParticipantTimelineData // Inhibitor kills per minute timeline counts
	Lane                            string                  // Participant's lane (legal values: MID, MIDDLE, TOP, JUNGLE, BOT, BOTTOM)
	Role                            string                  // Participant's role (legal values: DUO, NONE, SOLO, DUO_CARRY, DUO_SUPPORT)
	TowerAssistsPerMinCounts        ParticipantTimelineData // Tower assists per minute timeline counts
	TowerKillsPerMinCounts          ParticipantTimelineData // Tower kills per minute timeline counts
	TowerKillsPerMinDeltas          ParticipantTimelineData // Tower kills per minute timeline data
	VilemawAssistsPerMinCounts      ParticipantTimelineData // Vilemaw assists per minute timeline counts
	VilemawKillsPerMinCounts        ParticipantTimelineData // Vilemaw kills per minute timeline counts
	WardsPerMinDeltas               ParticipantTimelineData // Wards placed per minute timeline data
	XpDiffPerMinDeltas              ParticipantTimelineData // Experience difference per minute timeline data
	XpPerMinDeltas                  ParticipantTimelineData // Experience per minute timeline data
}

type Rune struct {
	RuneID int64 // Rune ID
	Rank   int64 // Rune rank
}

type Mastery struct {
	MasteryID int64 // Mastery ID
	Rank      int64 // Mastery rank
}

type BannedChampion struct {
	ChampionID int // Banned champion ID
	PickTurn   int // Turn during which the champion was banned
}

type Frame struct {
	Events            []Event                     // List of events for this frame.
	ParticipantFrames map[string]ParticipantFrame // Map of each participant ID to the participant's information for the frame.
	Timestamp         int64                       // Represents how many milliseconds into the game the frame occurred.
}

type ParticipantTimelineData struct {
	TenToTwenty    float64 // Value per minute from 10 min to 20 min
	ThirtyToEnd    float64 // Value per minute from 30 min to the end of the game
	TwentyToThirty float64 // Value per minute from 20 min to 30 min
	ZeroToTen      float64 // Value per minute from the beginning of the game to 10 min
}

type Event struct {
	// The ascended type of the event. Only present if relevant. Note that
	// CLEAR_ASCENDED refers to when a participants kills the ascended player.
	// (legal values: CHAMPION_ASCENDED, CLEAR_ASCENDED, MINION_ASCENDED)
	AscendedType string
	// The assisting participant IDs of the event. Only present if relevant.
	AssistingParticipantIDs []int
	// The building type of the event. Only present if relevant. (legal values:
	// INHIBITOR_BUILDING, TOWER_BUILDING)
	BuildingType string
	// The creator ID of the event. Only present if relevant.
	CreatorID int
	// Event type. (legal values: ASCENDED_EVENT, BUILDING_KILL, CAPTURE_POINT,
	// CHAMPION_KILL, ELITE_MONSTER_KILL, ITEM_DESTROYED, ITEM_PURCHASED,
	// ITEM_SOLD, ITEM_UNDO, SKILL_LEVEL_UP, WARD_KILL, WARD_PLACED)
	EventType string
	// The ending item ID of the event. Only present if relevant.
	ItemAfter int
	// The starting item ID of the event. Only present if relevant.
	ItemBefore int
	// The item ID of the event. Only present if relevant.
	ItemID int
	// The killer ID of the event. Only present if relevant.
	KillerID int
	// The lane type of the event. Only present if relevant. (legal values:
	// BOT_LANE, MID_LANE, TOP_LANE)
	LaneType string
	// The level up type of the event. Only present if relevant. (legal values:
	// EVOLVE, NORMAL)
	LevelUpType string
	// The monster type of the event. Only present if relevant. (legal values:
	// BARON_NASHOR, BLUE_GOLEM, DRAGON, RED_LIZARD, VILEMAW)
	MonsterType string
	// The participant ID of the event. Only present if relevant.
	ParticipantID int
	// The point captured in the event. Only present if relevant. (legal values:
	// POINT_A, POINT_B, POINT_C, POINT_D, POINT_E)
	PointCaptured string
	// The position of the event. Only present if relevant.
	Position Position
	// The skill slot of the event. Only present if relevant.
	SkillSlot int
	// The team ID of the event. Only present if relevant.
	TeamID int
	// Represents how many milliseconds into the game the event occurred.
	Timestamp int64
	// The tower type of the event. Only present if relevant. (legal values:
	// BASE_TURRET, FOUNTAIN_TURRET, INNER_TURRET, NEXUS_TURRET, OUTER_TURRET,
	// UNDEFINED_TURRET)
	TowerType string
	// The victim ID of the event. Only present if relevant.
	VictimID int
	// The ward type of the event. Only present if relevant. (legal values:
	// SIGHT_WARD, TEEMO_MUSHROOM, UNDEFINED, VISION_WARD, YELLOW_TRINKET,
	// YELLOW_TRINKET_UPGRADE)
	WardType string
}

type ParticipantFrame struct {
	CurrentGold         int      // Participant's current gold
	JungleMinionsKilled int      // Number of jungle minions killed by participant
	Level               int      // Participant's current level
	MinionsKilled       int      // Number of minions killed by participant
	ParticipantID       int      // Participant ID
	Position            Position // Participant's position
	TotalGold           int      // Participant's total gold
	XpPerMinDeltas      int      // Experience earned by participant
}

type Position struct {
	X, Y int
}
