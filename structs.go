package main

type Summoner struct {
	ID            int64  // Summoner ID.
	Name          string // Summoner name.
	ProfileIconID int    // ID of the summoner icon associated with the summoner.
	RevisionDate  int64  // Date summoner was last modified specified as epoch milliseconds.
	SummonerLevel int64  // Summoner level associated with the summoner.
}

type RecentGames struct {
	Games      []Game // Collection of recent games played (max 10).
	SummonerID int64  // Summoner ID.
}

type Game struct {
	ChampionID    int      // Champion ID associated with game.
	CreateDate    int64    // Date that end game data was recorded, specified as epoch milliseconds.
	FellowPlayers []Player // Other players associated with the game.
	GameID        int64    // Game ID.
	// Game mode: CLASSIC, ODIN, ARAM, TUTORIAL, ONEFORALL, FIRSTBLOOD.
	GameMode string
	// Game type: CUSTOM_GAME, MATCHED_GAME, TUTORIAL_GAME.
	GameType string
	Invalid  bool     // Invalid flag.
	IpEarned int      // IP Earned.
	Level    int      // Level.
	MapID    int      // Map ID.
	Spell1   int      // ID of first summoner spell.
	Spell2   int      // ID of second summoner spell.
	Stats    RawStats // Statistics associated with the game for this summoner.
	/* Game sub-type: NONE, NORMAL, BOT, RANKED_SOLO_5x5, RANKED_PREMADE_3x3,
	RANKED_PREMADE_5x5, ODIN_UNRANKED, RANKED_TEAM_3x3, RANKED_TEAM_5x5,
	NORMAL_3x3, BOT_3x3, CAP_5x5, ARAM_UNRANKED_5x5, ONEFORALL_5x5,
	FIRSTBLOOD_1x1, FIRSTBLOOD_2x2, SR_6x6, URF, URF_BOT, NIGHTMARE_BOT. */
	SubType string
	TeamID  int // Team ID associated with game. Team ID 100 is blue team. Team ID 200 is purple team.
}

type Player struct {
	ChampionID int   // ChampionintChampion id associated with player.
	SummonerID int64 // Summoner id associated with player.
	TeamID     int   // TeamintTeam id associated with player.
}

type RawStats struct {
	Assists                         int
	BarracksKilled                  int // Number of enemy inhibitors killed.
	ChampionsKilled                 int
	CombatPlayerScore               int
	ConsumablesPurchased            int
	DamageDealtPlayer               int
	DoubleKills                     int
	FirstBlood                      int
	Gold                            int
	GoldEarned                      int
	GoldSpent                       int
	Item0                           int
	Item1                           int
	Item2                           int
	Item3                           int
	Item4                           int
	Item5                           int
	Item6                           int
	ItemsPurchased                  int
	KillingSprees                   int
	LargestCriticalStrike           int
	LargestKillingSpree             int
	LargestMultiKill                int
	LegendaryItemsCreated           int // Number of tier 3 items built.
	Level                           int
	MagicDamageDealtPlayer          int
	MagicDamageDealtToChampions     int
	MagicDamageTaken                int
	MinionsDenied                   int
	MinionsKilled                   int
	NeutralMinionsKilled            int
	NeutralMinionsKilledEnemyJungle int
	NeutralMinionsKilledYourJungle  int
	NexusKilled                     bool // Flag specifying if the summoner got the killing blow on the nexus.
	NodeCapture                     int
	NodeCaptureAssist               int
	NodeNeutralize                  int
	NodeNeutralizeAssist            int
	NumDeaths                       int
	NumItemsBought                  int
	ObjectivePlayerScore            int
	PentaKills                      int
	PhysicalDamageDealtPlayer       int
	PhysicalDamageDealtToChampions  int
	PhysicalDamageTaken             int
	QuadraKills                     int
	SightWardsBought                int
	Spell1Cast                      int // Number of times first champion spell was cast.
	Spell2Cast                      int // Number of times second champion spell was cast.
	Spell3Cast                      int // Number of times third champion spell was cast.
	Spell4Cast                      int // Number of times fourth champion spell was cast.
	SummonSpell1Cast                int
	SummonSpell2Cast                int
	SuperMonsterKilled              int
	Team                            int
	TeamObjective                   int
	TimePlayed                      int
	TotalDamageDealt                int
	TotalDamageDealtToChampions     int
	TotalDamageTaken                int
	TotalHeal                       int
	TotalPlayerScore                int
	TotalScoreRank                  int
	TotalTimeCrowdControlDealt      int
	TotalUnitsHealed                int
	TripleKills                     int
	TrueDamageDealtPlayer           int
	TrueDamageDealtToChampions      int
	TrueDamageTaken                 int
	TurretsKilled                   int
	UnrealKills                     int
	VictoryPointTotal               int
	VisionWardsBought               int
	WardKilled                      int
	WardPlaced                      int
	Win                             bool // Flag specifying whether or not this game was won.
}

type RankedStats struct {
	ChampionsKilled []ChampionStats // Collection of aggregated stats summarized by champion.
	ModifyDate      int64           // Date stats were last modified specified as epoch milliseconds.
	SummonerID      int64           // Summoner ID.
}

type ChampionStats struct {
	/* Champion ID. Note that champion ID 0 represents the combined stats for all
	champions. For static information correlating to champion IDs, please refer to
	the LoL Static Data API. */
	ID    int
	Stats AggregatedStats // Aggregated stats associated with the champion.
}

type PlayerStatsSummaryList struct {
	PlayerStatSummaries []PlayerStatsSummary //Collection of player stats summaries associated with the summoner.
	SummonerID          int64                // Summoner ID.
}

type PlayerStatsSummary struct {
	AggregatedStats AggregatedStats // Aggregated stats.
	Losses          int             // Number of losses for this queue type. Returned for ranked queue types only.
	ModifyDate      int64           // Date stats were last modified specified as epoch milliseconds.
	/* Player stats summary type: AramUnranked5x5, CoopVsAI, CoopVsAI3x3,
	OdinUnranked, RankedPremade3x3, RankedPremade5x5, RankedSolo5x5,
	RankedTeam3x3, RankedTeam5x5, Unranked, Unranked3x3, OneForAll5x5,
	FirstBlood1x1, FirstBlood2x2, SummonersRift6x6, CAP5x5, URF, URFBots,
	NightmareBot. */
	PlayerStatSummaryType string
	Wins                  int // Number of wins for this queue type.
}

type AggregatedStats struct {
	AverageAssists              int // Dominion only.
	AverageChampionsKilled      int // Dominion only.
	AverageCombatPlayerScore    int // Dominion only.
	AverageNodeCapture          int // Dominion only.
	AverageNodeCaptureAssist    int // Dominion only.
	AverageNodeNeutralize       int // Dominion only.
	AverageNodeNeutralizeAssist int // Dominion only.
	AverageNumDeaths            int // Dominion only.
	AverageObjectivePlayerScore int // Dominion only.
	AverageTeamObjective        int // Dominion only.
	AverageTotalPlayerScore     int // Dominion only.
	BotGamesPlayed              int
	KillingSpree                int
	MaxAssists                  int // Dominion only.
	MaxChampionsKilled          int
	MaxCombatPlayerScore        int // Dominion only.
	MaxLargestCriticalStrike    int
	MaxLargestKillingSpree      int
	MaxNodeCapture              int // Dominion only.
	MaxNodeCaptureAssist        int // Dominion only.
	MaxNodeNeutralize           int // Dominion only.
	MaxNodeNeutralizeAssist     int // Dominion only.
	MaxNumDeaths                int // Only returned for ranked statistics.
	MaxObjectivePlayerScore     int // Dominion only.
	MaxTeamObjective            int // Dominion only.
	MaxTimePlayed               int
	MaxTimeSpentLiving          int
	MaxTotalPlayerScore         int // Dominion only.
	MostChampionKillsPerSession int
	MostSpellsCast              int
	NormalGamesPlayed           int
	RankedPremadeGamesPlayed    int
	RankedSoloGamesPlayed       int
	TotalAssists                int
	TotalChampionKills          int
	TotalDamageDealt            int
	TotalDamageTaken            int
	TotalDeathsPerSession       int // Only returned for ranked statistics.
	TotalDoubleKills            int
	TotalFirstBlood             int
	TotalGoldEarned             int
	TotalHeal                   int
	TotalMagicDamageDealt       int
	TotalMinionKills            int
	TotalNeutralMinionsKilled   int
	TotalNodeCapture            int // Dominion only.
	TotalNodeNeutralize         int // Dominion only.
	TotalPentaKills             int
	TotalPhysicalDamageDealt    int
	TotalQuadraKills            int
	TotalSessionsLost           int
	TotalSessionsPlayed         int
	TotalSessionsWon            int
	TotalTripleKills            int
	TotalTurretsKilled          int
	TotalUnrealKills            int
}
