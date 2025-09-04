// API Types based on OpenAPI specification

export interface IdeaInput {
  title: string;
  one_liner: string;
  category?: string;
  location?: string;
}

export interface ApproxLocation {
  country?: string;
  region?: string;
}

export interface AnalysisOptions {
  max_evidence?: number;
  location?: ApproxLocation;
  timeout?: string;
}

export interface AnalysisRequest {
  idea: IdeaInput;
  options?: AnalysisOptions;
}

export interface AnalysisResponse {
  analysis_id: string;
  status: 'completed' | 'failed';
}

export interface Evidence {
  id: string;
  url: string;
  title: string;
  snippet?: string;
  published_at?: string;
  retrieved_at: string;
  source_type?: string;
}

export interface Competitor {
  name: string;
  description: string;
  funding?: string;
  stage?: string;
  evidence_ids: string[];
}

export interface Risk {
  category: string;
  description: string;
  severity: number;
  likelihood: number;
  mitigation?: string;
  evidence_ids: string[];
}

export interface Barrier {
  type: 'regulation' | 'supply' | 'distribution' | 'trust' | 'tech';
  description: string;
  weight: number;
  evidence_ids: string[];
}

export interface GraveyardCase {
  company_name: string;
  description: string;
  failure_cause: string;
  lessons: string;
  evidence_ids: string[];
}

export interface MarketAnalysis {
  competitors: Competitor[];
  market_stage: 'early' | 'growing' | 'mature' | 'declining';
  positioning: string;
  evidence_ids: string[];
}

export interface ProblemAnalysis {
  pain_points: string[];
  validation: string;
  evidence_ids: string[];
}

export interface BarrierAnalysis {
  barriers: Barrier[];
  evidence_ids: string[];
}

export interface ExecutionAnalysis {
  capital_requirement: 'low' | 'medium' | 'high';
  talent_rarity: 'common' | 'rare' | 'very_rare';
  integration_count: number;
  complexity: number;
  evidence_ids: string[];
}

export interface RiskAnalysis {
  risks: Risk[];
  evidence_ids: string[];
}

export interface GraveyardAnalysis {
  cases: GraveyardCase[];
  evidence_ids: string[];
}

export interface Viability {
  overall_score: number;
  market_score: number;
  problem_score: number;
  barrier_score: number;
  execution_score: number;
  risk_score: number;
  graveyard_score: number;
  recommendation: string;
  key_insights: string[];
  evidence_ids: string[];
}

export interface Analysis {
  id: string;
  idea: IdeaInput;
  market: MarketAnalysis;
  problem: ProblemAnalysis;
  barriers: BarrierAnalysis;
  execution: ExecutionAnalysis;
  risks: RiskAnalysis;
  graveyard: GraveyardAnalysis;
  verdict: Viability;
  evidence: Evidence[];
  created_at: string;
  partial?: boolean;
  meta?: any;
}

export interface AnalysisListResponse {
  analyses: Analysis[];
  pagination: {
    limit: number;
    offset: number;
    total: number;
  };
}

export interface StatsResponse {
  total_analyses: number;
  max_evidence: number;
  timeout: string;
}

export interface HealthResponse {
  status: 'healthy' | 'unhealthy';
}

export interface ErrorResponse {
  error: string;
  code?: string;
  details?: string;
}

