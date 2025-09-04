import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
    Container,
    Typography,
    Box,
    Card,
    CardContent,
    Grid,
    Chip,
    CircularProgress,
    Alert,
    Button,
    Divider,
    List,
    ListItem,
    ListItemText,
    Accordion,
    AccordionSummary,
    AccordionDetails,
    Paper,
    Fade,
} from '@mui/material';
import {
    ExpandMore,
    TrendingUp,
    Psychology,
    Build,
    Warning,
    Security,
    Delete,
    Download,
    Share,
} from '@mui/icons-material';
import { useTheme } from '@mui/material/styles';
import { apiService } from '../services/apiService';
import { type Analysis } from '../types/api';
import ScoreCard from '../components/ScoreCard';

export default function AnalysisPage() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const theme = useTheme();
    const [analysis, setAnalysis] = useState<Analysis | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (!id) {
            navigate('/');
            return;
        }

        const fetchAnalysis = async () => {
            try {
                setLoading(true);
                setError(null);
                const result = await apiService.getAnalysis(id);
                setAnalysis(result);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Failed to load analysis');
            } finally {
                setLoading(false);
            }
        };

        fetchAnalysis();
    }, [id, navigate]);

    const handleDownloadMarkdown = async () => {
        if (!id) return;
        try {
            const markdown = await apiService.getAnalysisMarkdown(id);
            const blob = new Blob([markdown], { type: 'text/markdown' });
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `${analysis?.idea.title || 'analysis'}.md`;
            a.click();
            URL.revokeObjectURL(url);
        } catch (err) {
            console.error('Failed to download markdown:', err);
        }
    };

    const getRecommendationColor = (recommendation: string) => {
        if (recommendation.toLowerCase().includes('go')) return theme.palette.success.main;
        if (recommendation.toLowerCase().includes('no-go')) return theme.palette.error.main;
        return theme.palette.warning.main;
    };

    if (loading) {
        return (
            <Container maxWidth="lg" sx={{ py: 8, textAlign: 'center' }}>
                <CircularProgress size={60} />
                <Typography variant="h6" sx={{ mt: 2 }}>
                    Loading analysis...
                </Typography>
            </Container>
        );
    }

    if (error || !analysis) {
        return (
            <Container maxWidth="lg" sx={{ py: 8 }}>
                <Alert severity="error" sx={{ mb: 3 }}>
                    {error || 'Analysis not found'}
                </Alert>
                <Button variant="contained" onClick={() => navigate('/')}>
                    Go Back
                </Button>
            </Container>
        );
    }

    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            <Fade in timeout={800}>
                <Box>
                    {/* Header */}
                    <Box sx={{ mb: 4 }}>
                        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
                            <Box>
                                <Typography
                                    variant="h3"
                                    sx={{
                                        fontWeight: 700,
                                        background: theme.palette.gradient.primary,
                                        backgroundClip: 'text',
                                        WebkitBackgroundClip: 'text',
                                        WebkitTextFillColor: 'transparent',
                                        mb: 1,
                                    }}
                                >
                                    {analysis.idea.title}
                                </Typography>
                                <Typography variant="h6" color="text.secondary" sx={{ mb: 2 }}>
                                    {analysis.idea.one_liner}
                                </Typography>
                                {analysis.idea.category && (
                                    <Chip label={analysis.idea.category} sx={{ mr: 1 }} />
                                )}
                                {analysis.idea.location && (
                                    <Chip label={analysis.idea.location} variant="outlined" />
                                )}
                                {analysis.partial && (
                                    <Chip
                                        label="Partial Analysis"
                                        color="warning"
                                        sx={{ ml: 1 }}
                                    />
                                )}
                            </Box>

                            <Box sx={{ display: 'flex', gap: 1 }}>
                                <Button
                                    variant="outlined"
                                    startIcon={<Download />}
                                    onClick={handleDownloadMarkdown}
                                >
                                    Export
                                </Button>
                                <Button
                                    variant="outlined"
                                    startIcon={<Share />}
                                    onClick={() => navigator.clipboard.writeText(window.location.href)}
                                >
                                    Share
                                </Button>
                            </Box>
                        </Box>

                        <Typography variant="body2" color="text.secondary">
                            Analysis completed on {new Date(analysis.created_at).toLocaleDateString()}
                        </Typography>
                    </Box>

                    {/* Overall Verdict */}
                    <Paper
                        elevation={3}
                        sx={{
                            p: 4,
                            mb: 4,
                            background: `linear-gradient(135deg, rgba(255, 255, 255, 0.9) 0%, rgba(248, 250, 252, 0.9) 100%)`,
                            border: `2px solid ${getRecommendationColor(analysis.verdict.recommendation)}`,
                            borderRadius: 3,
                        }}
                    >
                        <Grid container spacing={3} alignItems="center">
                            <Grid item xs={12} md={8}>
                                <Typography variant="h4" sx={{ fontWeight: 700, mb: 1 }}>
                                    Overall Score: {analysis.verdict.overall_score.toFixed(1)}/100
                                </Typography>
                                <Typography
                                    variant="h5"
                                    sx={{
                                        color: getRecommendationColor(analysis.verdict.recommendation),
                                        fontWeight: 600,
                                        mb: 2,
                                    }}
                                >
                                    {analysis.verdict.recommendation}
                                </Typography>
                                <Typography variant="body1" sx={{ mb: 2 }}>
                                    Key Insights:
                                </Typography>
                                <List dense>
                                    {analysis.verdict.key_insights.map((insight, index) => (
                                        <ListItem key={index} sx={{ px: 0 }}>
                                            <ListItemText primary={`• ${insight}`} />
                                        </ListItem>
                                    ))}
                                </List>
                            </Grid>
                            <Grid item xs={12} md={4} sx={{ textAlign: 'center' }}>
                                <Box
                                    sx={{
                                        width: 120,
                                        height: 120,
                                        borderRadius: '50%',
                                        background: `conic-gradient(${getRecommendationColor(analysis.verdict.recommendation)} ${analysis.verdict.overall_score * 3.6}deg, ${theme.palette.grey[200]} 0deg)`,
                                        display: 'flex',
                                        alignItems: 'center',
                                        justifyContent: 'center',
                                        mx: 'auto',
                                        position: 'relative',
                                    }}
                                >
                                    <Box
                                        sx={{
                                            width: 80,
                                            height: 80,
                                            borderRadius: '50%',
                                            backgroundColor: 'white',
                                            display: 'flex',
                                            alignItems: 'center',
                                            justifyContent: 'center',
                                        }}
                                    >
                                        <Typography variant="h4" sx={{ fontWeight: 700 }}>
                                            {Math.round(analysis.verdict.overall_score)}
                                        </Typography>
                                    </Box>
                                </Box>
                            </Grid>
                        </Grid>
                    </Paper>

                    {/* Score Breakdown */}
                    <Typography variant="h4" sx={{ mb: 3, fontWeight: 600 }}>
                        Dimension Scores
                    </Typography>
                    <Grid container spacing={3} sx={{ mb: 4 }}>
                        <Grid item xs={12} sm={6} md={4}>
                            <ScoreCard
                                title="Market"
                                score={analysis.verdict.market_score}
                                description="Competition and positioning analysis"
                            />
                        </Grid>
                        <Grid item xs={12} sm={6} md={4}>
                            <ScoreCard
                                title="Problem"
                                score={analysis.verdict.problem_score}
                                description="Pain point validation and evidence"
                            />
                        </Grid>
                        <Grid item xs={12} sm={6} md={4}>
                            <ScoreCard
                                title="Barriers"
                                score={analysis.verdict.barrier_score}
                                description="Execution barriers and challenges"
                            />
                        </Grid>
                        <Grid item xs={12} sm={6} md={4}>
                            <ScoreCard
                                title="Execution"
                                score={analysis.verdict.execution_score}
                                description="Feasibility and resource requirements"
                            />
                        </Grid>
                        <Grid item xs={12} sm={6} md={4}>
                            <ScoreCard
                                title="Risks"
                                score={analysis.verdict.risk_score}
                                description="Business and market risks"
                            />
                        </Grid>
                        <Grid item xs={12} sm={6} md={4}>
                            <ScoreCard
                                title="Graveyard"
                                score={analysis.verdict.graveyard_score}
                                description="Lessons from failed competitors"
                            />
                        </Grid>
                    </Grid>

                    {/* Detailed Analysis */}
                    <Typography variant="h4" sx={{ mb: 3, fontWeight: 600 }}>
                        Detailed Analysis
                    </Typography>

                    <Grid container spacing={3}>
                        {/* Market Analysis */}
                        <Grid item xs={12}>
                            <Accordion defaultExpanded>
                                <AccordionSummary expandIcon={<ExpandMore />}>
                                    <TrendingUp sx={{ mr: 2 }} />
                                    <Typography variant="h6" sx={{ fontWeight: 600 }}>
                                        Market Analysis
                                    </Typography>
                                </AccordionSummary>
                                <AccordionDetails>
                                    <Grid container spacing={2}>
                                        <Grid item xs={12} md={6}>
                                            <Typography variant="subtitle1" sx={{ fontWeight: 600, mb: 1 }}>
                                                Market Stage: {analysis.market.market_stage}
                                            </Typography>
                                            <Typography variant="body2" sx={{ mb: 2 }}>
                                                {analysis.market.positioning}
                                            </Typography>
                                        </Grid>
                                        {analysis.market.competitors &&
                                            <Grid item xs={12} md={6}>
                                                <Typography variant="subtitle1" sx={{ fontWeight: 600, mb: 1 }}>
                                                    Competitors ({analysis.market.competitors?.length})
                                                </Typography>
                                                {analysis.market.competitors.slice(0, 3).map((competitor, index) => (
                                                    <Box key={index} sx={{ mb: 1 }}>
                                                        <Typography variant="body2" sx={{ fontWeight: 500 }}>
                                                            {competitor.name}
                                                        </Typography>
                                                        <Typography variant="caption" color="text.secondary">
                                                            {competitor.description}
                                                        </Typography>
                                                    </Box>
                                                ))}
                                            </Grid>
                                        }
                                    </Grid>
                                </AccordionDetails>
                            </Accordion>
                        </Grid>

                        {/* Problem Analysis */}
                        <Grid item xs={12}>
                            <Accordion>
                                <AccordionSummary expandIcon={<ExpandMore />}>
                                    <Psychology sx={{ mr: 2 }} />
                                    <Typography variant="h6" sx={{ fontWeight: 600 }}>
                                        Problem Analysis
                                    </Typography>
                                </AccordionSummary>
                                <AccordionDetails>
                                    <Typography variant="body2" sx={{ mb: 2 }}>
                                        {analysis.problem.validation}
                                    </Typography>
                                    {analysis.problem.pain_points.length > 0 && (
                                        <>
                                            <Typography variant="subtitle1" sx={{ fontWeight: 600, mb: 1 }}>
                                                Identified Pain Points:
                                            </Typography>
                                            <List dense>
                                                {analysis.problem.pain_points.map((point, index) => (
                                                    <ListItem key={index} sx={{ px: 0 }}>
                                                        <ListItemText primary={`• ${point}`} />
                                                    </ListItem>
                                                ))}
                                            </List>
                                        </>
                                    )}
                                </AccordionDetails>
                            </Accordion>
                        </Grid>

                        {/* Execution Analysis */}
                        <Grid item xs={12}>
                            <Accordion>
                                <AccordionSummary expandIcon={<ExpandMore />}>
                                    <Build sx={{ mr: 2 }} />
                                    <Typography variant="h6" sx={{ fontWeight: 600 }}>
                                        Execution Analysis
                                    </Typography>
                                </AccordionSummary>
                                <AccordionDetails>
                                    <Grid container spacing={2}>
                                        <Grid item xs={6} sm={3}>
                                            <Typography variant="subtitle2" color="text.secondary">
                                                Capital Requirement
                                            </Typography>
                                            <Chip label={analysis.execution.capital_requirement} size="small" />
                                        </Grid>
                                        <Grid item xs={6} sm={3}>
                                            <Typography variant="subtitle2" color="text.secondary">
                                                Talent Rarity
                                            </Typography>
                                            <Chip label={analysis.execution.talent_rarity} size="small" />
                                        </Grid>
                                        <Grid item xs={6} sm={3}>
                                            <Typography variant="subtitle2" color="text.secondary">
                                                Integrations
                                            </Typography>
                                            <Typography variant="body1">{analysis.execution.integration_count}</Typography>
                                        </Grid>
                                        <Grid item xs={6} sm={3}>
                                            <Typography variant="subtitle2" color="text.secondary">
                                                Complexity Score
                                            </Typography>
                                            <Typography variant="body1">{(analysis.execution.complexity * 100).toFixed(0)}%</Typography>
                                        </Grid>
                                    </Grid>
                                </AccordionDetails>
                            </Accordion>
                        </Grid>

                        {analysis.risks.risks && analysis.risks.risks.length > 0 && (
                            <Grid item xs={12}>
                                <Accordion>
                                    <AccordionSummary expandIcon={<ExpandMore />}>
                                        <Warning sx={{ mr: 2 }} />
                                        <Typography variant="h6" sx={{ fontWeight: 600 }}>
                                            Risk Analysis ({analysis.risks.risks.length} risks)
                                        </Typography>
                                    </AccordionSummary>
                                    <AccordionDetails>
                                        <Grid container spacing={2}>
                                            {analysis.risks.risks.map((risk, index) => (
                                                <Grid item xs={12} md={6} key={index}>
                                                    <Card variant="outlined" sx={{ p: 2 }}>
                                                        <Typography variant="subtitle1" sx={{ fontWeight: 600, mb: 1 }}>
                                                            {risk.category}
                                                        </Typography>
                                                        <Typography variant="body2" sx={{ mb: 1 }}>
                                                            {risk.description}
                                                        </Typography>
                                                        <Box sx={{ display: 'flex', gap: 1, mb: 1 }}>
                                                            <Chip
                                                                label={`Severity: ${risk.severity}/5`}
                                                                size="small"
                                                                color="error"
                                                            />
                                                            <Chip
                                                                label={`Likelihood: ${risk.likelihood}/5`}
                                                                size="small"
                                                                color="warning"
                                                            />
                                                        </Box>
                                                        {risk.mitigation && (
                                                            <Typography variant="caption" color="text.secondary">
                                                                Mitigation: {risk.mitigation}
                                                            </Typography>
                                                        )}
                                                    </Card>
                                                </Grid>
                                            ))}
                                        </Grid>
                                    </AccordionDetails>
                                </Accordion>
                            </Grid>
                        )}

                        {/* Graveyard Cases */}
                        {analysis.graveyard.cases.length > 0 && (
                            <Grid item xs={12}>
                                <Accordion>
                                    <AccordionSummary expandIcon={<ExpandMore />}>
                                        <Delete sx={{ mr: 2 }} />
                                        <Typography variant="h6" sx={{ fontWeight: 600 }}>
                                            Graveyard Analysis ({analysis.graveyard.cases.length} cases)
                                        </Typography>
                                    </AccordionSummary>
                                    <AccordionDetails>
                                        <Grid container spacing={2}>
                                            {analysis.graveyard.cases.map((graveyardCase, index) => (
                                                <Grid item xs={12} key={index}>
                                                    <Card variant="outlined" sx={{ p: 2 }}>
                                                        <Typography variant="subtitle1" sx={{ fontWeight: 600, mb: 1 }}>
                                                            {graveyardCase.company_name}
                                                        </Typography>
                                                        <Typography variant="body2" sx={{ mb: 1 }}>
                                                            {graveyardCase.description}
                                                        </Typography>
                                                        <Typography variant="body2" color="error.main" sx={{ mb: 1 }}>
                                                            <strong>Failure Cause:</strong> {graveyardCase.failure_cause}
                                                        </Typography>
                                                        <Typography variant="body2" color="text.secondary">
                                                            <strong>Lessons:</strong> {graveyardCase.lessons}
                                                        </Typography>
                                                    </Card>
                                                </Grid>
                                            ))}
                                        </Grid>
                                    </AccordionDetails>
                                </Accordion>
                            </Grid>
                        )}
                    </Grid>

                    {analysis.evidence && analysis.evidence.length > 0 && (
                        <Box sx={{ mt: 4 }}>
                            <Typography variant="h4" sx={{ mb: 2, fontWeight: 600 }}>
                                Evidence Summary
                            </Typography>
                            <Typography variant="body1" color="text.secondary">
                                This analysis is based on {analysis.evidence.length} pieces of evidence
                                from various sources including news articles, databases, forums, and academic sources.
                            </Typography>
                        </Box>
                    )}
                </Box>
            </Fade>
        </Container>
    );
}

