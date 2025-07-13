package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"url-shortener/internal/api/middleware"
	"url-shortener/internal/core/ports"
)

type CompetitiveIntelligenceHandler struct {
	competitiveService ports.CompetitiveIntelligenceService
}

func NewCompetitiveIntelligenceHandler(competitiveService ports.CompetitiveIntelligenceService) *CompetitiveIntelligenceHandler {
	return &CompetitiveIntelligenceHandler{
		competitiveService: competitiveService,
	}
}

// Market Analysis

// GetMarketPosition handles getting user's market position analysis
func (h *CompetitiveIntelligenceHandler) GetMarketPosition(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	position, err := h.competitiveService.AnalyzeMarketPosition(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, "Failed to analyze market position", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, position, http.StatusOK)
}

// GetCompetitorData handles getting specific competitor information
func (h *CompetitiveIntelligenceHandler) GetCompetitorData(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	competitorID := chi.URLParam(r, "competitorId")
	if competitorID == "" {
		h.writeErrorResponse(w, "Competitor ID is required", http.StatusBadRequest)
		return
	}

	competitor, err := h.competitiveService.GetCompetitorData(r.Context(), competitorID)
	if err != nil {
		h.writeErrorResponse(w, "Failed to get competitor data", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, competitor, http.StatusOK)
}

// GetMarketTrends handles getting market trends for an industry
func (h *CompetitiveIntelligenceHandler) GetMarketTrends(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	industry := r.URL.Query().Get("industry")
	if industry == "" {
		industry = "URL Shortening" // Default industry
	}

	trends, err := h.competitiveService.GetMarketTrends(r.Context(), industry)
	if err != nil {
		h.writeErrorResponse(w, "Failed to get market trends", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, trends, http.StatusOK)
}

// Benchmarking

// GetIndustryBenchmarks handles getting industry benchmark data
func (h *CompetitiveIntelligenceHandler) GetIndustryBenchmarks(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	industry := r.URL.Query().Get("industry")
	if industry == "" {
		industry = "URL Shortening" // Default industry
	}

	benchmarks, err := h.competitiveService.GetIndustryBenchmarks(r.Context(), industry)
	if err != nil {
		h.writeErrorResponse(w, "Failed to get industry benchmarks", http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, benchmarks, http.StatusOK)
}

// ComparePerformance handles comparing user performance against a competitor
func (h *CompetitiveIntelligenceHandler) ComparePerformance(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	competitorID := chi.URLParam(r, "competitorId")
	if competitorID == "" {
		h.writeErrorResponse(w, "Competitor ID is required", http.StatusBadRequest)
		return
	}

	comparison, err := h.competitiveService.ComparePerformance(r.Context(), userID, competitorID)
	if err != nil {
		h.writeErrorResponse(w, "Failed to compare performance", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"user_id":      userID,
		"competitor_id": competitorID,
		"comparison":   comparison,
		"summary": map[string]interface{}{
			"overall_performance": comparison["overall_performance"],
			"performance_rating": h.getPerformanceRating(comparison["overall_performance"]),
			"top_strengths":      h.getTopStrengths(comparison),
			"improvement_areas":  h.getImprovementAreas(comparison),
		},
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// GetPerformanceGaps handles getting performance gaps and improvement opportunities
func (h *CompetitiveIntelligenceHandler) GetPerformanceGaps(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	gaps, err := h.competitiveService.GetPerformanceGaps(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, "Failed to get performance gaps", http.StatusInternalServerError)
		return
	}

	// Calculate summary metrics
	totalGaps := len(gaps)
	highPriorityGaps := 0
	totalPotentialImpact := 0.0

	for _, gap := range gaps {
		if gap.Priority == "high" {
			highPriorityGaps++
		}
		totalPotentialImpact += gap.GapPercentage
	}

	response := map[string]interface{}{
		"gaps": gaps,
		"summary": map[string]interface{}{
			"total_gaps":             totalGaps,
			"high_priority_gaps":     highPriorityGaps,
			"total_potential_impact": totalPotentialImpact,
			"avg_gap_percentage":     func() float64 {
				if totalGaps > 0 {
					return totalPotentialImpact / float64(totalGaps)
				}
				return 0
			}(),
		},
		"recommendations": h.getPriorityRecommendations(gaps),
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// Opportunity Identification

// GetMarketOpportunities handles getting market opportunities based on trends
func (h *CompetitiveIntelligenceHandler) GetMarketOpportunities(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	opportunities, err := h.competitiveService.IdentifyMarketOpportunities(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, "Failed to identify market opportunities", http.StatusInternalServerError)
		return
	}

	// Filter and prioritize opportunities
	prioritizedOpportunities := h.prioritizeOpportunities(opportunities)

	response := map[string]interface{}{
		"opportunities": prioritizedOpportunities,
		"summary": map[string]interface{}{
			"total_opportunities": len(opportunities),
			"high_impact":         h.countByImpact(opportunities, "high"),
			"medium_impact":       h.countByImpact(opportunities, "medium"),
			"low_impact":          h.countByImpact(opportunities, "low"),
		},
		"strategic_focus": h.getStrategicFocus(prioritizedOpportunities),
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// GetCompetitorWeaknesses handles getting weaknesses of a specific competitor
func (h *CompetitiveIntelligenceHandler) GetCompetitorWeaknesses(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	competitorID := chi.URLParam(r, "competitorId")
	if competitorID == "" {
		h.writeErrorResponse(w, "Competitor ID is required", http.StatusBadRequest)
		return
	}

	weaknesses, err := h.competitiveService.AnalyzeCompetitorWeaknesses(r.Context(), competitorID)
	if err != nil {
		h.writeErrorResponse(w, "Failed to analyze competitor weaknesses", http.StatusInternalServerError)
		return
	}

	// Get competitor data for context
	competitor, err := h.competitiveService.GetCompetitorData(r.Context(), competitorID)
	if err != nil {
		h.writeErrorResponse(w, "Failed to get competitor data", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"competitor_id":   competitorID,
		"competitor_name": competitor.Name,
		"weaknesses":      weaknesses,
		"opportunities":   h.convertWeaknessesToOpportunities(weaknesses),
		"attack_vectors":  h.getAttackVectors(weaknesses),
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// GetEmergingTrends handles getting emerging trends in the industry
func (h *CompetitiveIntelligenceHandler) GetEmergingTrends(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	industry := r.URL.Query().Get("industry")
	if industry == "" {
		industry = "URL Shortening" // Default industry
	}

	trends, err := h.competitiveService.GetEmergingTrends(r.Context(), industry)
	if err != nil {
		h.writeErrorResponse(w, "Failed to get emerging trends", http.StatusInternalServerError)
		return
	}

	// Filter trends by confidence and impact
	minConfidence := 70.0 // Only show trends with >70% confidence
	filteredTrends := []interface{}{}
	
	for _, trend := range trends {
		if trend.Confidence >= minConfidence {
			trendData := map[string]interface{}{
				"name":        trend.Name,
				"impact":      trend.Impact,
				"confidence":  trend.Confidence,
				"description": trend.Description,
				"timeline":    trend.Timeline,
				"adoption":    trend.Adoption,
				"opportunity_score": h.calculateOpportunityScore(trend),
				"recommendations": h.getTrendActionPlan(trend.Name),
			}
			filteredTrends = append(filteredTrends, trendData)
		}
	}

	response := map[string]interface{}{
		"industry":        industry,
		"trends":          filteredTrends,
		"trend_count":     len(filteredTrends),
		"action_priority": h.getTrendPriorities(trends),
		"investment_recommendations": h.getInvestmentRecommendations(trends),
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// Advanced Analytics

// GetCompetitiveOverview handles getting a comprehensive competitive overview
func (h *CompetitiveIntelligenceHandler) GetCompetitiveOverview(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		h.writeErrorResponse(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Get market position
	position, err := h.competitiveService.AnalyzeMarketPosition(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, "Failed to get market position", http.StatusInternalServerError)
		return
	}

	// Get performance gaps
	gaps, err := h.competitiveService.GetPerformanceGaps(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, "Failed to get performance gaps", http.StatusInternalServerError)
		return
	}

	// Get market opportunities
	opportunities, err := h.competitiveService.IdentifyMarketOpportunities(r.Context(), userID)
	if err != nil {
		h.writeErrorResponse(w, "Failed to get market opportunities", http.StatusInternalServerError)
		return
	}

	// Get industry benchmarks
	benchmarks, err := h.competitiveService.GetIndustryBenchmarks(r.Context(), "URL Shortening")
	if err != nil {
		h.writeErrorResponse(w, "Failed to get benchmarks", http.StatusInternalServerError)
		return
	}

	overview := map[string]interface{}{
		"market_position": position,
		"performance_gaps": map[string]interface{}{
			"gaps":          gaps[:min(len(gaps), 5)], // Top 5 gaps
			"priority_count": h.countByPriority(gaps, "high"),
		},
		"opportunities": map[string]interface{}{
			"market_opportunities": opportunities[:min(len(opportunities), 3)], // Top 3 opportunities
			"quick_wins":          h.getQuickWins(gaps),
		},
		"benchmarks": map[string]interface{}{
			"key_metrics": h.getKeyBenchmarkMetrics(benchmarks),
			"your_tier":   h.getUserTier(position, benchmarks),
		},
		"strategic_recommendations": h.getStrategicRecommendations(position, gaps, opportunities),
		"competitive_score":         h.calculateCompetitiveScore(position, gaps),
	}

	h.writeJSONResponse(w, overview, http.StatusOK)
}

// Helper methods

func (h *CompetitiveIntelligenceHandler) getPerformanceRating(score float64) string {
	if score >= 120 {
		return "Excellent"
	} else if score >= 100 {
		return "Good"
	} else if score >= 80 {
		return "Average"
	}
	return "Needs Improvement"
}

func (h *CompetitiveIntelligenceHandler) getTopStrengths(comparison map[string]float64) []string {
	strengths := []string{}
	if comparison["click_through_rate_ratio"] > 110 {
		strengths = append(strengths, "Superior click-through rate")
	}
	if comparison["conversion_rate_ratio"] > 110 {
		strengths = append(strengths, "Higher conversion rate")
	}
	if comparison["volume_ratio"] > 110 {
		strengths = append(strengths, "Higher traffic volume")
	}
	return strengths
}

func (h *CompetitiveIntelligenceHandler) getImprovementAreas(comparison map[string]float64) []string {
	areas := []string{}
	if comparison["click_through_rate_ratio"] < 90 {
		areas = append(areas, "Click-through rate optimization")
	}
	if comparison["conversion_rate_ratio"] < 90 {
		areas = append(areas, "Conversion optimization")
	}
	if comparison["geographic_reach_ratio"] < 90 {
		areas = append(areas, "Geographic expansion")
	}
	return areas
}

func (h *CompetitiveIntelligenceHandler) getPriorityRecommendations(gaps []interface{}) []string {
	recommendations := []string{}
	// Add logic to extract high-priority recommendations
	recommendations = append(recommendations, "Focus on top 3 performance gaps")
	recommendations = append(recommendations, "Implement quick-win optimizations first")
	recommendations = append(recommendations, "Consider competitive positioning strategy")
	return recommendations
}

func (h *CompetitiveIntelligenceHandler) prioritizeOpportunities(opportunities []interface{}) []interface{} {
	// Simple prioritization - in a real implementation you'd sort by impact and feasibility
	return opportunities
}

func (h *CompetitiveIntelligenceHandler) countByImpact(opportunities []interface{}, impact string) int {
	count := 0
	// Add counting logic based on impact level
	return count
}

func (h *CompetitiveIntelligenceHandler) getStrategicFocus(opportunities []interface{}) []string {
	return []string{
		"AI-powered analytics implementation",
		"Privacy-first tracking solutions", 
		"Mobile optimization initiatives",
	}
}

func (h *CompetitiveIntelligenceHandler) convertWeaknessesToOpportunities(weaknesses []string) []string {
	opportunities := []string{}
	for _, weakness := range weaknesses {
		opportunity := "Capitalize on competitor's " + weakness
		opportunities = append(opportunities, opportunity)
	}
	return opportunities
}

func (h *CompetitiveIntelligenceHandler) getAttackVectors(weaknesses []string) []string {
	vectors := []string{}
	if len(weaknesses) > 0 {
		vectors = append(vectors, "Target their weak customer segments")
		vectors = append(vectors, "Highlight superior features in marketing")
		vectors = append(vectors, "Offer competitive switching incentives")
	}
	return vectors
}

func (h *CompetitiveIntelligenceHandler) calculateOpportunityScore(trend interface{}) float64 {
	// Simple scoring algorithm
	return 75.5 // Mock score
}

func (h *CompetitiveIntelligenceHandler) getTrendActionPlan(trendName string) []string {
	plans := map[string][]string{
		"AI-Powered Link Analytics": {
			"Research ML/AI integration options",
			"Develop proof of concept",
			"Plan full implementation roadmap",
		},
		"Privacy-First Tracking": {
			"Audit current privacy compliance",
			"Implement GDPR-friendly tracking",
			"Communicate privacy benefits to users",
		},
	}
	if plan, exists := plans[trendName]; exists {
		return plan
	}
	return []string{"Research trend implementation", "Develop strategy"}
}

func (h *CompetitiveIntelligenceHandler) getTrendPriorities(trends []interface{}) []map[string]interface{} {
	return []map[string]interface{}{
		{"trend": "AI-Powered Analytics", "priority": "high", "timeline": "6 months"},
		{"trend": "Privacy-First Tracking", "priority": "high", "timeline": "3 months"},
		{"trend": "Real-Time Collaboration", "priority": "medium", "timeline": "9 months"},
	}
}

func (h *CompetitiveIntelligenceHandler) getInvestmentRecommendations(trends []interface{}) []map[string]interface{} {
	return []map[string]interface{}{
		{"area": "AI/ML capabilities", "investment": "high", "roi_potential": "very_high"},
		{"area": "Privacy compliance", "investment": "medium", "roi_potential": "high"},
		{"area": "Mobile optimization", "investment": "medium", "roi_potential": "medium"},
	}
}

func (h *CompetitiveIntelligenceHandler) countByPriority(gaps []interface{}, priority string) int {
	// Add counting logic
	return 2 // Mock count
}

func (h *CompetitiveIntelligenceHandler) getQuickWins(gaps []interface{}) []map[string]interface{} {
	return []map[string]interface{}{
		{"area": "Mobile responsiveness", "effort": "low", "impact": "medium"},
		{"area": "Page load optimization", "effort": "medium", "impact": "high"},
	}
}

func (h *CompetitiveIntelligenceHandler) getKeyBenchmarkMetrics(benchmarks interface{}) map[string]interface{} {
	return map[string]interface{}{
		"industry_avg_ctr":        4.1,
		"industry_avg_conversion": 3.5,
		"industry_avg_mobile":     69.2,
	}
}

func (h *CompetitiveIntelligenceHandler) getUserTier(position interface{}, benchmarks interface{}) string {
	// Add logic to determine user's performance tier
	return "above_average"
}

func (h *CompetitiveIntelligenceHandler) getStrategicRecommendations(position, gaps, opportunities interface{}) []string {
	return []string{
		"Focus on conversion rate optimization for immediate impact",
		"Invest in AI-powered analytics for competitive advantage",
		"Expand geographic reach to capture new markets",
		"Implement privacy-first tracking to build trust",
	}
}

func (h *CompetitiveIntelligenceHandler) calculateCompetitiveScore(position, gaps interface{}) float64 {
	// Mock competitive score calculation
	return 78.5
}

// Helper functions

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (h *CompetitiveIntelligenceHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to encode response"}`))
	}
}

func (h *CompetitiveIntelligenceHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.Write([]byte(`{"error": "Internal server error"}`))
	}
}