import React, { useState, useEffect } from 'react';
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
    List,
    ListItem,
    ListItemText,
    ListItemIcon,
    Paper,
    Divider,
    Fade,
} from '@mui/material';
import {
    Analytics,
    Speed,
    Schedule,
    CheckCircle,
    TrendingUp,
    Psychology,
    Build,
    Warning,
    Security,
    Delete,
} from '@mui/icons-material';
import { useTheme } from '@mui/material/styles';
import { apiService } from '../services/apiService';
import type { StatsResponse, HealthResponse } from '../types/api';

export default function AboutPage() {
    const [stats, setStats] = useState<StatsResponse | null>(null);
    const [health, setHealth] = useState<HealthResponse | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const theme = useTheme();

    useEffect(() => {
        const fetchData = async () => {
            try {
                setLoading(true);
                setError(null);
                const [statsResponse, healthResponse] = await Promise.all([
                    apiService.getStats(),
                    apiService.healthCheck(),
                ]);
                setStats(statsResponse);
                setHealth(healthResponse);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Failed to load system information');
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, []);

    const analysisFramework = [
        {
            icon: <TrendingUp />,
            title: 'Market Analysis',
            weight: '25%',
            description: 'Competition landscape, market stage, and positioning opportunities',
            details: ['Competitor identification', 'Market maturity assessment', 'Positioning strategy'],
        },
        {
            icon: <Psychology />,
            title: 'Problem Validation',
            weight: '20%',
            description: 'Pain point validation and user evidence gathering',
            details: ['Pain point identification', 'Problem urgency', 'User evidence validation'],
        },
        {
            icon: <Security />,
            title: 'Barrier Analysis',
            weight: '15%',
            description: 'Execution barriers including regulatory and technical challenges',
            details: ['Regulatory hurdles', 'Technical challenges', 'Distribution difficulties'],
        },
        {
            icon: <Build />,
            title: 'Execution Feasibility',
            weight: '15%',
            description: 'Resource requirements and implementation complexity',
            details: ['Capital requirements', 'Talent availability', 'Complexity assessment'],
        },
        {
            icon: <Warning />,
            title: 'Risk Assessment',
            weight: '15%',
            description: 'Business, market, and technical risk evaluation',
            details: ['Business risks', 'Market risks', 'Technical risks'],
        },
        {
            icon: <Delete />,
            title: 'Graveyard Study',
            weight: '10%',
            description: 'Learning from failed competitors and similar ventures',
            details: ['Failed competitor analysis', 'Lesson extraction', 'Pattern recognition'],
        },
    ];

    if (loading) {
        return (
            <Container maxWidth="lg" sx={{ py: 8, textAlign: 'center' }}>
                <CircularProgress size={60} />
                <Typography variant="h6" sx={{ mt: 2 }}>
                    Loading system information...
                </Typography>
            </Container>
        );
    }

    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            <Fade in timeout={800}>
                <Box>
                    {/* Header */}
                    <Box sx={{ textAlign: 'center', mb: 6 }}>
                        <Typography
                            variant="h3"
                            sx={{
                                fontWeight: 700,
                                background: theme.palette.gradient.primary,
                                backgroundClip: 'text',
                                WebkitBackgroundClip: 'text',
                                WebkitTextFillColor: 'transparent',
                                mb: 2,
                            }}
                        >
                            About RealityCheck
                        </Typography>
                        <Typography variant="h6" color="text.secondary" sx={{ maxWidth: '600px', mx: 'auto' }}>
                            AI-powered startup idea validation through comprehensive multi-dimensional analysis
                        </Typography>
                    </Box>

                    {/* Error State */}
                    {error && (
                        <Alert severity="error" sx={{ mb: 3 }}>
                            {error}
                        </Alert>
                    )}

                    {/* System Status */}
                    <Paper
                        elevation={3}
                        sx={{
                            p: 4,
                            mb: 6,
                            background: 'rgba(255, 255, 255, 0.9)',
                            backdropFilter: 'blur(10px)',
                        }}
                    >
                        <Typography variant="h4" sx={{ mb: 3, fontWeight: 600 }}>
                            System Status
                        </Typography>
                        <Grid container spacing={3}>
                            <Grid item xs={12} sm={6} md={3}>
                                <Box sx={{ textAlign: 'center' }}>
                                    <CheckCircle
                                        sx={{
                                            fontSize: 48,
                                            color: health?.status === 'healthy' ? theme.palette.success.main : theme.palette.error.main,
                                            mb: 1,
                                        }}
                                    />
                                    <Typography variant="h6" sx={{ fontWeight: 600 }}>
                                        {health?.status === 'healthy' ? 'Healthy' : 'Unhealthy'}
                                    </Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        System Status
                                    </Typography>
                                </Box>
                            </Grid>
                            {stats && (
                                <>
                                    <Grid item xs={12} sm={6} md={3}>
                                        <Box sx={{ textAlign: 'center' }}>
                                            <Analytics sx={{ fontSize: 48, color: theme.palette.primary.main, mb: 1 }} />
                                            <Typography variant="h6" sx={{ fontWeight: 600 }}>
                                                {stats.total_analyses.toLocaleString()}
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                Total Analyses
                                            </Typography>
                                        </Box>
                                    </Grid>
                                    <Grid item xs={12} sm={6} md={3}>
                                        <Box sx={{ textAlign: 'center' }}>
                                            <Speed sx={{ fontSize: 48, color: theme.palette.secondary.main, mb: 1 }} />
                                            <Typography variant="h6" sx={{ fontWeight: 600 }}>
                                                {stats.max_evidence}
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                Max Evidence/Analysis
                                            </Typography>
                                        </Box>
                                    </Grid>
                                    <Grid item xs={12} sm={6} md={3}>
                                        <Box sx={{ textAlign: 'center' }}>
                                            <Schedule sx={{ fontSize: 48, color: theme.palette.accent.primary, mb: 1 }} />
                                            <Typography variant="h6" sx={{ fontWeight: 600 }}>
                                                {stats.timeout}
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                Analysis Timeout
                                            </Typography>
                                        </Box>
                                    </Grid>
                                </>
                            )}
                        </Grid>
                    </Paper>

                    {/* Analysis Framework */}
                    <Typography variant="h4" sx={{ mb: 4, fontWeight: 600, textAlign: 'center' }}>
                        Analysis Framework
                    </Typography>
                    <Grid container spacing={3} sx={{ mb: 6 }}>
                        {analysisFramework.map((dimension, index) => (
                            <Grid item xs={12} md={6} key={index}>
                                <Card
                                    elevation={2}
                                    sx={{
                                        height: '100%',
                                        background: 'rgba(255, 255, 255, 0.9)',
                                        backdropFilter: 'blur(10px)',
                                        border: '1px solid rgba(255, 255, 255, 0.2)',
                                        transition: 'all 0.3s ease-in-out',
                                        '&:hover': {
                                            transform: 'translateY(-4px)',
                                            boxShadow: theme.shadows[8],
                                        },
                                    }}
                                >
                                    <CardContent sx={{ p: 3 }}>
                                        <Box sx={{ display: 'flex', alignItems: 'flex-start', mb: 2 }}>
                                            <Box
                                                sx={{
                                                    p: 1,
                                                    borderRadius: 2,
                                                    background: theme.palette.gradient.primary,
                                                    color: 'white',
                                                    mr: 2,
                                                }}
                                            >
                                                {dimension.icon}
                                            </Box>
                                            <Box sx={{ flex: 1 }}>
                                                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
                                                    <Typography variant="h6" sx={{ fontWeight: 600 }}>
                                                        {dimension.title}
                                                    </Typography>
                                                    <Chip label={dimension.weight} size="small" color="primary" />
                                                </Box>
                                                <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                                                    {dimension.description}
                                                </Typography>
                                                <List dense>
                                                    {dimension.details.map((detail, detailIndex) => (
                                                        <ListItem key={detailIndex} sx={{ px: 0, py: 0.5 }}>
                                                            <ListItemText
                                                                primary={`• ${detail}`}
                                                                primaryTypographyProps={{ variant: 'body2' }}
                                                            />
                                                        </ListItem>
                                                    ))}
                                                </List>
                                            </Box>
                                        </Box>
                                    </CardContent>
                                </Card>
                            </Grid>
                        ))}
                    </Grid>

                    {/* How It Works */}
                    <Paper
                        elevation={2}
                        sx={{
                            p: 4,
                            background: 'rgba(255, 255, 255, 0.9)',
                            backdropFilter: 'blur(10px)',
                        }}
                    >
                        <Typography variant="h4" sx={{ mb: 3, fontWeight: 600, textAlign: 'center' }}>
                            How It Works
                        </Typography>
                        <Grid container spacing={4}>
                            <Grid item xs={12} md={4}>
                                <Box sx={{ textAlign: 'center' }}>
                                    <Box
                                        sx={{
                                            width: 80,
                                            height: 80,
                                            borderRadius: '50%',
                                            background: theme.palette.gradient.primary,
                                            color: 'white',
                                            display: 'flex',
                                            alignItems: 'center',
                                            justifyContent: 'center',
                                            mx: 'auto',
                                            mb: 2,
                                        }}
                                    >
                                        <Typography variant="h4" sx={{ fontWeight: 700 }}>
                                            1
                                        </Typography>
                                    </Box>
                                    <Typography variant="h5" sx={{ fontWeight: 600, mb: 1 }}>
                                        Submit Your Idea
                                    </Typography>
                                    <Typography variant="body1" color="text.secondary">
                                        Provide a title and description of your startup idea. Include optional details like category and location for better analysis.
                                    </Typography>
                                </Box>
                            </Grid>
                            <Grid item xs={12} md={4}>
                                <Box sx={{ textAlign: 'center' }}>
                                    <Box
                                        sx={{
                                            width: 80,
                                            height: 80,
                                            borderRadius: '50%',
                                            background: theme.palette.gradient.secondary,
                                            color: 'white',
                                            display: 'flex',
                                            alignItems: 'center',
                                            justifyContent: 'center',
                                            mx: 'auto',
                                            mb: 2,
                                        }}
                                    >
                                        <Typography variant="h4" sx={{ fontWeight: 700 }}>
                                            2
                                        </Typography>
                                    </Box>
                                    <Typography variant="h5" sx={{ fontWeight: 600, mb: 1 }}>
                                        AI Analysis
                                    </Typography>
                                    <Typography variant="body1" color="text.secondary">
                                        Our AI system conducts comprehensive research across multiple dimensions, gathering evidence from various sources and analyzing market dynamics.
                                    </Typography>
                                </Box>
                            </Grid>
                            <Grid item xs={12} md={4}>
                                <Box sx={{ textAlign: 'center' }}>
                                    <Box
                                        sx={{
                                            width: 80,
                                            height: 80,
                                            borderRadius: '50%',
                                            background: theme.palette.gradient.tertiary,
                                            color: 'white',
                                            display: 'flex',
                                            alignItems: 'center',
                                            justifyContent: 'center',
                                            mx: 'auto',
                                            mb: 2,
                                        }}
                                    >
                                        <Typography variant="h4" sx={{ fontWeight: 700 }}>
                                            3
                                        </Typography>
                                    </Box>
                                    <Typography variant="h5" sx={{ fontWeight: 600, mb: 1 }}>
                                        Get Insights
                                    </Typography>
                                    <Typography variant="body1" color="text.secondary">
                                        Receive a detailed report with scores, recommendations, and actionable insights to help you make informed decisions about your startup idea.
                                    </Typography>
                                </Box>
                            </Grid>
                        </Grid>
                    </Paper>
                </Box>
            </Fade>
        </Container>
    );
}

