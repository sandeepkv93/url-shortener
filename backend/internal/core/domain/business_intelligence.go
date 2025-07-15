package domain

import (
	"time"

	"gorm.io/gorm"
)

// Business Intelligence Dashboard Models

// DashboardWidget represents a configurable widget in the BI dashboard
type DashboardWidget struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	Type        string         `json:"type" gorm:"not null"` // chart, metric, table, map, funnel
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	Position    WidgetPosition `json:"position" gorm:"embedded"`
	Size        WidgetSize     `json:"size" gorm:"embedded"`
	Config      WidgetConfig   `json:"config" gorm:"type:jsonb"`
	IsVisible   bool           `json:"is_visible" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type WidgetPosition struct {
	X int `json:"x" gorm:"not null"`
	Y int `json:"y" gorm:"not null"`
	Z int `json:"z" gorm:"default:0"`
}

type WidgetSize struct {
	Width  int `json:"width" gorm:"not null"`
	Height int `json:"height" gorm:"not null"`
}

type WidgetConfig struct {
	ChartType     string                 `json:"chart_type,omitempty"`     // line, bar, pie, donut, area
	DataSource    string                 `json:"data_source,omitempty"`    // clicks, urls, devices, geo
	Filters       map[string]interface{} `json:"filters,omitempty"`        // date range, url filters
	AggregationType string               `json:"aggregation_type,omitempty"` // sum, avg, count, unique
	GroupBy       string                 `json:"group_by,omitempty"`       // date, country, device
	TimeRange     string                 `json:"time_range,omitempty"`     // 1h, 24h, 7d, 30d, 90d, 1y
	RefreshRate   int                    `json:"refresh_rate,omitempty"`   // seconds
	Colors        []string               `json:"colors,omitempty"`         // custom color scheme
	ShowLegend    bool                   `json:"show_legend,omitempty"`
	ShowGrid      bool                   `json:"show_grid,omitempty"`
	Limit         int                    `json:"limit,omitempty"`          // max results
}

// Dashboard represents a custom dashboard layout
type Dashboard struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	IsDefault   bool           `json:"is_default" gorm:"default:false"`
	IsPublic    bool           `json:"is_public" gorm:"default:false"`
	Layout      DashboardLayout `json:"layout" gorm:"type:jsonb"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User    *User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Widgets []DashboardWidget `json:"widgets,omitempty" gorm:"foreignKey:UserID"`
}

type DashboardLayout struct {
	GridSize    GridSize   `json:"grid_size"`
	Breakpoints Breakpoint `json:"breakpoints"`
	Margin      []int      `json:"margin"`     // [horizontal, vertical]
	Padding     []int      `json:"padding"`    // [top, right, bottom, left]
	RowHeight   int        `json:"row_height"` // pixels
}

type GridSize struct {
	Columns int `json:"columns"`
	Rows    int `json:"rows"`
}

type Breakpoint struct {
	XS int `json:"xs"` // extra small screens
	SM int `json:"sm"` // small screens
	MD int `json:"md"` // medium screens
	LG int `json:"lg"` // large screens
	XL int `json:"xl"` // extra large screens
}

// BI Analytics Models

// FunnelStep represents a step in a conversion funnel
type FunnelStep struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	FunnelID    uint           `json:"funnel_id" gorm:"not null;index"`
	StepNumber  int            `json:"step_number" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	URLPattern  string         `json:"url_pattern"` // regex pattern for matching URLs
	EventType   string         `json:"event_type"`  // click, view, conversion
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Funnel *ConversionFunnel `json:"funnel,omitempty" gorm:"foreignKey:FunnelID"`
}

// ConversionFunnel represents a user-defined conversion funnel
type ConversionFunnel struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User  *User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Steps []FunnelStep `json:"steps,omitempty" gorm:"foreignKey:FunnelID"`
}

// FunnelAnalytics represents funnel conversion analytics
type FunnelAnalytics struct {
	FunnelID        uint                   `json:"funnel_id"`
	TotalEntries    int64                  `json:"total_entries"`
	Conversions     int64                  `json:"conversions"`
	ConversionRate  float64                `json:"conversion_rate"`
	DropOffRate     float64                `json:"drop_off_rate"`
	StepAnalytics   []FunnelStepAnalytics  `json:"step_analytics"`
	TimeToConvert   map[string]interface{} `json:"time_to_convert"`
	ConversionTrend map[string]int64       `json:"conversion_trend"`
}

type FunnelStepAnalytics struct {
	StepNumber     int     `json:"step_number"`
	StepName       string  `json:"step_name"`
	Entries        int64   `json:"entries"`
	Exits          int64   `json:"exits"`
	Conversions    int64   `json:"conversions"`
	ConversionRate float64 `json:"conversion_rate"`
	DropOffRate    float64 `json:"drop_off_rate"`
	AvgTimeOnStep  float64 `json:"avg_time_on_step"` // seconds
}

// AdvancedAnalytics represents comprehensive business intelligence data
type AdvancedAnalytics struct {
	UserID              uint                    `json:"user_id"`
	Period              string                  `json:"period"`
	GeneratedAt         time.Time               `json:"generated_at"`
	PerformanceMetrics  BIPerformanceMetrics    `json:"performance_metrics"`
	AudienceInsights    AudienceInsights        `json:"audience_insights"`
	ContentAnalytics    ContentAnalytics        `json:"content_analytics"`
	CompetitiveAnalysis CompetitiveAnalysis     `json:"competitive_analysis"`
	PredictiveInsights  PredictiveInsights      `json:"predictive_insights"`
	RecommendationEngine RecommendationEngine  `json:"recommendation_engine"`
}

type BIPerformanceMetrics struct {
	CTR                 float64               `json:"ctr"`                   // Click-through rate
	ConversionRate      float64               `json:"conversion_rate"`
	BounceRate          float64               `json:"bounce_rate"`
	EngagementScore     float64               `json:"engagement_score"`
	QualityScore        float64               `json:"quality_score"`
	PerformanceTrend    map[string]float64    `json:"performance_trend"`
	BenchmarkComparison map[string]interface{} `json:"benchmark_comparison"`
}

type AudienceInsights struct {
	Demographics     Demographics         `json:"demographics"`
	BehaviorPatterns BehaviorPatterns     `json:"behavior_patterns"`
	SegmentAnalysis  []AudienceSegment    `json:"segment_analysis"`
	RetentionMetrics RetentionMetrics     `json:"retention_metrics"`
	UserJourney      []UserJourneyStep    `json:"user_journey"`
}

type Demographics struct {
	AgeGroups     map[string]int64 `json:"age_groups"`
	Gender        map[string]int64 `json:"gender"`
	Locations     map[string]int64 `json:"locations"`
	Devices       map[string]int64 `json:"devices"`
	Languages     map[string]int64 `json:"languages"`
	Interests     map[string]int64 `json:"interests"`
}

type BehaviorPatterns struct {
	ClickPatterns    map[string]interface{} `json:"click_patterns"`
	TimeOfDay        map[int]int64          `json:"time_of_day"`
	DayOfWeek        map[string]int64       `json:"day_of_week"`
	SeasonalTrends   map[string]int64       `json:"seasonal_trends"`
	ReturnVisitors   int64                  `json:"return_visitors"`
	AverageSession   float64                `json:"average_session"`
}

type AudienceSegment struct {
	SegmentID       string                 `json:"segment_id"`
	Name            string                 `json:"name"`
	Size            int64                  `json:"size"`
	Characteristics map[string]interface{} `json:"characteristics"`
	Performance     map[string]float64     `json:"performance"`
	GrowthRate      float64                `json:"growth_rate"`
}

type RetentionMetrics struct {
	RetentionRate      float64            `json:"retention_rate"`
	ChurnRate          float64            `json:"churn_rate"`
	LifetimeValue      float64            `json:"lifetime_value"`
	RetentionCohort    map[string]float64 `json:"retention_cohort"`
	ReactivationRate   float64            `json:"reactivation_rate"`
}

type UserJourneyStep struct {
	StepNumber      int                    `json:"step_number"`
	Action          string                 `json:"action"`
	URL             string                 `json:"url"`
	Timestamp       time.Time              `json:"timestamp"`
	Duration        float64                `json:"duration"` // seconds
	InteractionData map[string]interface{} `json:"interaction_data"`
}

type ContentAnalytics struct {
	TopPerformingContent []ContentPerformance `json:"top_performing_content"`
	ContentCategories    map[string]int64     `json:"content_categories"`
	EngagementMetrics    ContentEngagement    `json:"engagement_metrics"`
	ContentLifecycle     ContentLifecycle     `json:"content_lifecycle"`
	ViralityMetrics      ViralityMetrics      `json:"virality_metrics"`
}

type ContentPerformance struct {
	ContentID       string  `json:"content_id"`
	Title           string  `json:"title"`
	URL             string  `json:"url"`
	Clicks          int64   `json:"clicks"`
	Shares          int64   `json:"shares"`
	EngagementScore float64 `json:"engagement_score"`
	Category        string  `json:"category"`
	Tags            []string `json:"tags"`
}

type ContentEngagement struct {
	AverageTimeOnPage float64            `json:"average_time_on_page"`
	ShareRate         float64            `json:"share_rate"`
	CommentRate       float64            `json:"comment_rate"`
	LikeRate          float64            `json:"like_rate"`
	EngagementTrend   map[string]float64 `json:"engagement_trend"`
}

type ContentLifecycle struct {
	CreationDate    time.Time `json:"creation_date"`
	PeakDate        time.Time `json:"peak_date"`
	DeclineDate     time.Time `json:"decline_date"`
	LifespanDays    int       `json:"lifespan_days"`
	PeakTraffic     int64     `json:"peak_traffic"`
	CurrentStatus   string    `json:"current_status"` // growing, stable, declining
}

type ViralityMetrics struct {
	ViralCoefficient float64            `json:"viral_coefficient"`
	ShareVelocity    float64            `json:"share_velocity"`
	ReachAmplification float64          `json:"reach_amplification"`
	ViralPaths       []ViralPath        `json:"viral_paths"`
}

type ViralPath struct {
	Source      string  `json:"source"`
	Destination string  `json:"destination"`
	Shares      int64   `json:"shares"`
	Strength    float64 `json:"strength"`
}

type CompetitiveAnalysis struct {
	MarketPosition   MarketPosition   `json:"market_position"`
	CompetitorMetrics []CompetitorData `json:"competitor_metrics"`
	BenchmarkData    BenchmarkData    `json:"benchmark_data"`
	MarketTrends     MarketTrends     `json:"market_trends"`
	OpportunityGaps  []OpportunityGap `json:"opportunity_gaps"`
}

type MarketPosition struct {
	UserID                uint                `json:"user_id"`
	Rank                  int                 `json:"rank"`
	Tier                  string              `json:"tier"`
	Percentile            float64             `json:"percentile"`
	MarketShare           float64             `json:"market_share"`
	Ranking               int                 `json:"ranking"`
	GrowthRate            float64             `json:"growth_rate"`
	CompetitiveGap        float64             `json:"competitive_gap"`
	Metrics               PositionMetrics     `json:"metrics"`
	Strengths             []string            `json:"strengths"`
	Weaknesses            []string            `json:"weaknesses"`
	Opportunities         []string            `json:"opportunities"`
	Threats               []string            `json:"threats"`
	CompetitiveAdvantages []string            `json:"competitive_advantages"`
	RecommendedActions    []string            `json:"recommended_actions"`
	UpdatedAt             time.Time           `json:"updated_at"`
}

type PositionMetrics struct {
	TotalClicks     int64   `json:"total_clicks"`
	ClicksPerURL    float64 `json:"clicks_per_url"`
	ConversionRate  float64 `json:"conversion_rate"`
	RetentionRate   float64 `json:"retention_rate"`
	GrowthRate      float64 `json:"growth_rate"`
	ClickShare      float64 `json:"click_share"`
	URLShare        float64 `json:"url_share"`
	EngagementScore float64 `json:"engagement_score"`
	ReachScore      float64 `json:"reach_score"`
}

type CompetitorData struct {
	ID              string              `json:"id"`
	CompetitorID    string              `json:"competitor_id"`
	Name            string              `json:"name"`
	Industry        string              `json:"industry"`
	Website         string              `json:"website"`
	Description     string              `json:"description"`
	MarketShare     float64             `json:"market_share"`
	Founded         int                 `json:"founded"`
	Employees       int                 `json:"employees"`
	Revenue         float64             `json:"revenue"`
	GrowthRate      float64             `json:"growth_rate"`
	PerformanceGap  float64             `json:"performance_gap"`
	KeyMetrics      map[string]interface{} `json:"key_metrics"`
	Metrics         CompetitorMetrics   `json:"metrics"`
	Strengths       []string            `json:"strengths"`
	Weaknesses      []string            `json:"weaknesses"`
	RecentNews      []NewsItem          `json:"recent_news"`
	SocialMedia     SocialMediaData     `json:"social_media"`
	TechnologyStack []string            `json:"technology_stack"`
	Partnerships    []string            `json:"partnerships"`
	LastUpdated     time.Time           `json:"last_updated"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

type CompetitorMetrics struct {
	MonthlyActiveUsers int64   `json:"monthly_active_users"`
	URLsCreated        int64   `json:"urls_created"`
	ClicksPerMonth     int64   `json:"clicks_per_month"`
	ConversionRate     float64 `json:"conversion_rate"`
	AverageClickTime   float64 `json:"average_click_time"`
	GeographicReach    float64 `json:"geographic_reach"`
	MobileUsage        float64 `json:"mobile_usage"`
}

type SocialMediaData struct {
	TwitterFollowers  int64 `json:"twitter_followers"`
	LinkedInFollowers int64 `json:"linkedin_followers"`
	FacebookFollowers int64 `json:"facebook_followers"`
	InstagramFollowers int64 `json:"instagram_followers"`
	YouTubeSubscribers int64 `json:"youtube_subscribers"`
}

type BenchmarkData struct {
	Industry        string              `json:"industry"`
	DataSource      string              `json:"data_source"`
	UpdatedAt       time.Time           `json:"updated_at"`
	Metrics         []BenchmarkMetric   `json:"metrics"`
	IndustryAverage map[string]float64  `json:"industry_average"`
	TopPerformers   map[string]float64  `json:"top_performers"`
	YourPerformance map[string]float64  `json:"your_performance"`
	Percentile      map[string]float64  `json:"percentile"`
}

type BenchmarkMetric struct {
	Name         string  `json:"name"`
	Value        float64 `json:"value"`
	Unit         string  `json:"unit"`
	Description  string  `json:"description"`
	Tier         PerformanceTier `json:"tier"`
	Median       float64 `json:"median"`
	Average      float64 `json:"average"`
	Percentile25 float64 `json:"percentile_25"`
	Percentile75 float64 `json:"percentile_75"`
	Percentile90 float64 `json:"percentile_90"`
}

type PerformanceTier string

const (
	PerformanceTierExcellent PerformanceTier = "excellent"
	PerformanceTierGood      PerformanceTier = "good"
	PerformanceTierAverage   PerformanceTier = "average"
	PerformanceTierPoor      PerformanceTier = "poor"
)

type MarketTrends struct {
	Industry              string             `json:"industry"`
	TimeRange             string             `json:"time_range"`
	UpdatedAt             time.Time          `json:"updated_at"`
	GrowthMetrics         GrowthMetrics      `json:"growth_metrics"`
	EmergingTrends        []Trend            `json:"emerging_trends"`
	SeasonalPatterns      map[string]float64 `json:"seasonal_patterns"`
	ConsumerBehavior      map[string]interface{} `json:"consumer_behavior"`
	TechnologyTrends      []string           `json:"technology_trends"`
	TechnologyAdoption    map[string]float64 `json:"technology_adoption"`
	CustomerBehavior      CustomerBehavior   `json:"customer_behavior"`
	CompetitiveLandscape  CompetitiveLandscape `json:"competitive_landscape"`
}

type CustomerBehavior struct {
	AverageSessionDuration          float64            `json:"average_session_duration"`
	BounceRate                      float64            `json:"bounce_rate"`
	ConversionRate                  float64            `json:"conversion_rate"`
	PreferredFeatures               []string           `json:"preferred_features"`
	UsagePatterns                   map[string]float64 `json:"usage_patterns"`
	MobileUsage                     float64            `json:"mobile_usage"`
	ReturnVisitorRate               float64            `json:"return_visitor_rate"`
	ConversionFunnelOptimization    float64            `json:"conversion_funnel_optimization"`
}

type CompetitiveLandscape struct {
	MarketLeaders     []string           `json:"market_leaders"`
	EmergingPlayers   []string           `json:"emerging_players"`
	MarketConcentration float64          `json:"market_concentration"`
	CompetitiveIntensity string          `json:"competitive_intensity"`
	BarriersToEntry   []string           `json:"barriers_to_entry"`
}

type GrowthMetrics struct {
	MarketSize                    float64 `json:"market_size"`
	YearOverYearGrowth           float64 `json:"year_over_year_growth"`
	ProjectedGrowth              float64 `json:"projected_growth"`
	CompoundAnnualGrowthRate     float64 `json:"compound_annual_growth_rate"`
	MarketSaturation             float64 `json:"market_saturation"`
	AdoptionRate                 float64 `json:"adoption_rate"`
}

type Trend struct {
	Name        string  `json:"name"`
	Direction   string  `json:"direction"` // up, down, stable
	Impact      string  `json:"impact"`    // high, medium, low
	Confidence  float64 `json:"confidence"`
	Description string  `json:"description"`
	Timeline    string  `json:"timeline"`
	Adoption    float64 `json:"adoption"`
}

type OpportunityGap struct {
	Area                 string   `json:"area"`
	Impact               string   `json:"impact"`
	Effort               string   `json:"effort"`
	Priority             string   `json:"priority"`
	Description          string   `json:"description"`
	Potential            float64  `json:"potential"`
	CurrentValue         float64  `json:"current_value"`
	BenchmarkValue       float64  `json:"benchmark_value"`
	GapPercentage        float64  `json:"gap_percentage"`
	Recommendations      []string `json:"recommendations"`
	EstimatedImpact      string   `json:"estimated_impact"`
	ImplementationEffort string   `json:"implementation_effort"`
	Timeline             string   `json:"timeline"`
}

type PredictiveInsights struct {
	ForecastData    ForecastData    `json:"forecast_data"`
	AnomalyDetection []Anomaly       `json:"anomaly_detection"`
	TrendPrediction TrendPrediction `json:"trend_prediction"`
	RiskAssessment  RiskAssessment  `json:"risk_assessment"`
	OptimizationSuggestions []OptimizationSuggestion `json:"optimization_suggestions"`
}

type ForecastData struct {
	Period      string           `json:"period"`
	Predictions map[string]float64 `json:"predictions"`
	Confidence  map[string]float64 `json:"confidence"`
	Scenarios   map[string]Scenario `json:"scenarios"` // best, worst, most likely
}

type Scenario struct {
	Name        string             `json:"name"`
	Probability float64            `json:"probability"`
	Predictions map[string]float64 `json:"predictions"`
	Impact      string             `json:"impact"`
}

type Anomaly struct {
	Type        string    `json:"type"`        // spike, drop, outlier
	Metric      string    `json:"metric"`
	Severity    string    `json:"severity"`    // high, medium, low
	Detected    time.Time `json:"detected"`
	Value       float64   `json:"value"`
	Expected    float64   `json:"expected"`
	Deviation   float64   `json:"deviation"`
	Description string    `json:"description"`
	Impact      string    `json:"impact"`
}

type TrendPrediction struct {
	ShortTerm  TrendForecast `json:"short_term"`  // next 7 days
	MediumTerm TrendForecast `json:"medium_term"` // next 30 days
	LongTerm   TrendForecast `json:"long_term"`   // next 90 days
}

type TrendForecast struct {
	Direction   string  `json:"direction"`
	Magnitude   float64 `json:"magnitude"`
	Confidence  float64 `json:"confidence"`
	FactorsInfluencing []string `json:"factors_influencing"`
}

type RiskAssessment struct {
	OverallRisk RiskLevel     `json:"overall_risk"`
	RiskFactors []RiskFactor  `json:"risk_factors"`
	Mitigation  []Mitigation  `json:"mitigation"`
}

type RiskLevel struct {
	Level       string  `json:"level"`       // low, medium, high, critical
	Score       float64 `json:"score"`       // 0-100
	Description string  `json:"description"`
	Trend       string  `json:"trend"`       // increasing, stable, decreasing
}

type RiskFactor struct {
	Factor      string  `json:"factor"`
	Impact      string  `json:"impact"`
	Probability float64 `json:"probability"`
	Description string  `json:"description"`
}

type Mitigation struct {
	Strategy    string  `json:"strategy"`
	Priority    string  `json:"priority"`
	Effort      string  `json:"effort"`
	Effectiveness float64 `json:"effectiveness"`
	Description string  `json:"description"`
}

type OptimizationSuggestion struct {
	Area        string  `json:"area"`
	Suggestion  string  `json:"suggestion"`
	Impact      string  `json:"impact"`
	Effort      string  `json:"effort"`
	Priority    string  `json:"priority"`
	Potential   float64 `json:"potential"`
	Confidence  float64 `json:"confidence"`
	Implementation string `json:"implementation"`
}

type RecommendationEngine struct {
	PersonalizedRecommendations []Recommendation `json:"personalized_recommendations"`
	ContentRecommendations      []ContentRec     `json:"content_recommendations"`
	OptimizationRecommendations []OptimizationRec `json:"optimization_recommendations"`
	AudienceRecommendations     []AudienceRec    `json:"audience_recommendations"`
}

type Recommendation struct {
	Type        string  `json:"type"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Priority    string  `json:"priority"`
	Impact      string  `json:"impact"`
	Confidence  float64 `json:"confidence"`
	ActionItems []string `json:"action_items"`
}

type ContentRec struct {
	ContentType string   `json:"content_type"`
	Topics      []string `json:"topics"`
	Timing      string   `json:"timing"`
	Channels    []string `json:"channels"`
	Potential   float64  `json:"potential"`
}

type OptimizationRec struct {
	Metric      string  `json:"metric"`
	CurrentValue float64 `json:"current_value"`
	TargetValue float64 `json:"target_value"`
	Strategy    string  `json:"strategy"`
	Timeline    string  `json:"timeline"`
}

type AudienceRec struct {
	SegmentName    string   `json:"segment_name"`
	Characteristics []string `json:"characteristics"`
	Size           int64    `json:"size"`
	Potential      float64  `json:"potential"`
	Approach       string   `json:"approach"`
}

// ScheduledReport represents automated reporting configuration
type ScheduledReport struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	UserID       uint           `json:"user_id" gorm:"not null;index"`
	Name         string         `json:"name" gorm:"not null"`
	Description  string         `json:"description"`
	ReportType   string         `json:"report_type"` // dashboard, analytics, funnel, competitive
	Schedule     string         `json:"schedule"`    // cron expression
	Recipients   []string       `json:"recipients" gorm:"type:jsonb"`
	Format       string         `json:"format"`      // pdf, excel, csv, json
	Config       ReportConfig   `json:"config" gorm:"type:jsonb"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	LastRun      *time.Time     `json:"last_run"`
	NextRun      *time.Time     `json:"next_run"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type ReportConfig struct {
	DateRange    string                 `json:"date_range"`    // 7d, 30d, 90d, custom
	CustomRange  DateRange              `json:"custom_range,omitempty"`
	Metrics      []string               `json:"metrics"`       // specific metrics to include
	Filters      map[string]interface{} `json:"filters"`       // data filters
	Grouping     string                 `json:"grouping"`      // daily, weekly, monthly
	Comparison   bool                   `json:"comparison"`    // include period comparison
	Charts       bool                   `json:"charts"`        // include visualizations
	RawData      bool                   `json:"raw_data"`      // include raw data export
}

// DataExport represents export request and metadata
type DataExport struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	ExportType  string         `json:"export_type"`  // analytics, users, urls, clicks
	Format      string         `json:"format"`       // csv, excel, json, pdf
	Status      string         `json:"status"`       // pending, processing, completed, failed
	FilePath    string         `json:"file_path"`
	FileSize    int64          `json:"file_size"`
	RecordCount int64          `json:"record_count"`
	Config      ExportConfig   `json:"config" gorm:"type:jsonb"`
	ExpiresAt   time.Time      `json:"expires_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type ExportConfig struct {
	DateRange   DateRange              `json:"date_range"`
	Filters     map[string]interface{} `json:"filters"`
	Columns     []string               `json:"columns"`      // specific columns to export
	Compression bool                   `json:"compression"`  // compress the export file
	Encryption  bool                   `json:"encryption"`   // encrypt the export file
}

// Request/Response DTOs

type CreateDashboardRequest struct {
	Name        string          `json:"name" validate:"required,max=100"`
	Description string          `json:"description" validate:"max=500"`
	IsDefault   bool            `json:"is_default"`
	IsPublic    bool            `json:"is_public"`
	Layout      DashboardLayout `json:"layout"`
}

type UpdateDashboardRequest struct {
	Name        *string          `json:"name,omitempty" validate:"omitempty,max=100"`
	Description *string          `json:"description,omitempty" validate:"omitempty,max=500"`
	IsDefault   *bool            `json:"is_default,omitempty"`
	IsPublic    *bool            `json:"is_public,omitempty"`
	Layout      *DashboardLayout `json:"layout,omitempty"`
}

type CreateWidgetRequest struct {
	Type        string         `json:"type" validate:"required,oneof=chart metric table map funnel"`
	Title       string         `json:"title" validate:"required,max=100"`
	Description string         `json:"description" validate:"max=500"`
	Position    WidgetPosition `json:"position" validate:"required"`
	Size        WidgetSize     `json:"size" validate:"required"`
	Config      WidgetConfig   `json:"config" validate:"required"`
}

type UpdateWidgetRequest struct {
	Title       *string         `json:"title,omitempty" validate:"omitempty,max=100"`
	Description *string         `json:"description,omitempty" validate:"omitempty,max=500"`
	Position    *WidgetPosition `json:"position,omitempty"`
	Size        *WidgetSize     `json:"size,omitempty"`
	Config      *WidgetConfig   `json:"config,omitempty"`
	IsVisible   *bool           `json:"is_visible,omitempty"`
}

type CreateFunnelRequest struct {
	Name        string        `json:"name" validate:"required,max=100"`
	Description string        `json:"description" validate:"max=500"`
	Steps       []FunnelStepRequest `json:"steps" validate:"required,min=2,max=10"`
}

type FunnelStepRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"max=500"`
	URLPattern  string `json:"url_pattern" validate:"required"`
	EventType   string `json:"event_type" validate:"required,oneof=click view conversion"`
}

type CreateScheduledReportRequest struct {
	Name        string       `json:"name" validate:"required,max=100"`
	Description string       `json:"description" validate:"max=500"`
	ReportType  string       `json:"report_type" validate:"required,oneof=dashboard analytics funnel competitive"`
	Schedule    string       `json:"schedule" validate:"required"`    // cron expression
	Recipients  []string     `json:"recipients" validate:"required,dive,email"`
	Format      string       `json:"format" validate:"required,oneof=pdf excel csv json"`
	Config      ReportConfig `json:"config" validate:"required"`
}

type DataExportRequest struct {
	ExportType string       `json:"export_type" validate:"required,oneof=analytics users urls clicks"`
	Format     string       `json:"format" validate:"required,oneof=csv excel json pdf"`
	Config     ExportConfig `json:"config" validate:"required"`
}

// Response DTOs

type DashboardResponse struct {
	*Dashboard
	WidgetCount int `json:"widget_count"`
}

type DashboardListResponse struct {
	Dashboards []DashboardResponse `json:"dashboards"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}

type WidgetDataResponse struct {
	WidgetID  uint                   `json:"widget_id"`
	Data      interface{}            `json:"data"`
	Metadata  map[string]interface{} `json:"metadata"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type FunnelAnalyticsResponse struct {
	*FunnelAnalytics
	RecommendedOptimizations []OptimizationSuggestion `json:"recommended_optimizations"`
}

// Validation methods

func (r *CreateDashboardRequest) Validate() error {
	if r.Name == "" {
		return ErrInvalidRequest
	}
	return nil
}

func (r *CreateWidgetRequest) Validate() error {
	if r.Type == "" || r.Title == "" {
		return ErrInvalidRequest
	}
	validTypes := []string{"chart", "metric", "table", "map", "funnel"}
	for _, validType := range validTypes {
		if r.Type == validType {
			return nil
		}
	}
	return ErrInvalidRequest
}

func (r *CreateFunnelRequest) Validate() error {
	if r.Name == "" || len(r.Steps) < 2 {
		return ErrInvalidRequest
	}
	return nil
}

func (r *DataExportRequest) Validate() error {
	if r.ExportType == "" || r.Format == "" {
		return ErrInvalidRequest
	}
	return nil
}

// Job status for async operations like report generation and data exports
type JobStatus struct {
	ID          string                 `json:"id" gorm:"primarykey"`
	Type        string                 `json:"type" gorm:"not null"` // report_generation, data_export, funnel_analysis
	Status      string                 `json:"status" gorm:"not null"` // pending, running, completed, failed
	Progress    float64                `json:"progress" gorm:"default:0"`
	Message     string                 `json:"message"`
	Result      map[string]interface{} `json:"result" gorm:"type:jsonb"`
	Error       string                 `json:"error"`
	UserID      uint                   `json:"user_id" gorm:"not null;index"`
	ResourceID  uint                   `json:"resource_id"` // report ID, export ID, etc.
	StartedAt   time.Time              `json:"started_at"`
	CompletedAt *time.Time             `json:"completed_at"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	
	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}


// Additional domain types for reporting and execution

type ReportExecution struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	ReportID    uint           `json:"report_id" gorm:"not null;index"`
	Status      string         `json:"status"` // running, completed, failed
	StartedAt   time.Time      `json:"started_at"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	Duration    *int64         `json:"duration,omitempty"` // milliseconds
	FilePath    string         `json:"file_path"`
	FileSize    int64          `json:"file_size"`
	ErrorMessage string        `json:"error_message"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Report *ScheduledReport `json:"report,omitempty" gorm:"foreignKey:ReportID"`
}