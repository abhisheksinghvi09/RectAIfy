import type {
    AnalysisRequest,
    AnalysisResponse,
    Analysis,
    AnalysisListResponse,
    StatsResponse,
    HealthResponse,
    ErrorResponse,
} from '../types/api';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:9444';

class ApiService {
    private async request<T>(
        endpoint: string,
        options: RequestInit = {}
    ): Promise<T> {
        const url = `${API_BASE_URL}${endpoint}`;

        const config: RequestInit = {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
            ...options,
        };

        // Add bearer token if available
        const token = localStorage.getItem('api_token');
        if (token) {
            config.headers = {
                ...config.headers,
                'Authorization': `Bearer ${token}`,
            };
        }

        try {
            const response = await fetch(url, config);

            if (!response.ok) {
                const errorData: ErrorResponse = await response.json().catch(() => ({
                    error: `HTTP ${response.status}: ${response.statusText}`,
                }));
                throw new Error(errorData.error || 'Request failed');
            }

            // Handle different content types
            const contentType = response.headers.get('content-type');
            if (contentType?.includes('application/json')) {
                return await response.json();
            } else if (contentType?.includes('text/')) {
                return await response.text() as unknown as T;
            }

            return await response.json();
        } catch (error) {
            if (error instanceof Error) {
                throw error;
            }
            throw new Error('Network error occurred');
        }
    }

    // Health check
    async healthCheck(): Promise<HealthResponse> {
        return this.request<HealthResponse>('/health');
    }

    // Submit idea for analysis
    async submitAnalysis(request: AnalysisRequest): Promise<AnalysisResponse> {
        return this.request<AnalysisResponse>('/v1/analyze', {
            method: 'POST',
            body: JSON.stringify(request),
        });
    }

    // Get analysis results
    async getAnalysis(id: string): Promise<Analysis> {
        return this.request<Analysis>(`/v1/analyses/${id}`);
    }

    // Get analysis as markdown
    async getAnalysisMarkdown(id: string): Promise<string> {
        return this.request<string>(`/v1/analyses/${id}.md`);
    }

    // Get analysis as HTML
    async getAnalysisHTML(id: string): Promise<string> {
        return this.request<string>(`/v1/analyses/${id}.html`);
    }

    // List analyses with pagination and search
    async listAnalyses(
        limit: number = 10,
        offset: number = 0,
        query?: string
    ): Promise<AnalysisListResponse> {
        const params = new URLSearchParams({
            limit: limit.toString(),
            offset: offset.toString(),
        });

        if (query) {
            params.append('q', query);
        }

        return this.request<AnalysisListResponse>(`/v1/analyses?${params}`);
    }

    // Get system statistics
    async getStats(): Promise<StatsResponse> {
        return this.request<StatsResponse>('/v1/stats');
    }
}

export const apiService = new ApiService();
export default apiService;

