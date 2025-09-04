import React, { useState, useEffect } from 'react';
import {
    Container,
    Typography,
    Box,
    Card,
    CardContent,
    Grid,
    TextField,
    InputAdornment,
    Chip,
    CircularProgress,
    Alert,
    Button,
    Pagination,
    Paper,
    IconButton,
    Tooltip,
    Fade,
} from '@mui/material';
import {
    Search,
    TrendingUp,
    AccessTime,
    Visibility,
    Delete,
    Refresh,
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { useTheme } from '@mui/material/styles';
import { apiService } from '../services/apiService';
import { type Analysis, type AnalysisListResponse } from '../types/api';

export default function DashboardPage() {
    const [analyses, setAnalyses] = useState<Analysis[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [searchQuery, setSearchQuery] = useState('');
    const [currentPage, setCurrentPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [totalCount, setTotalCount] = useState(0);
    const pageSize = 12;
    const navigate = useNavigate();
    const theme = useTheme();

    const fetchAnalyses = async (page: number = 1, query: string = '') => {
        try {
            setLoading(true);
            setError(null);
            const offset = (page - 1) * pageSize;
            const response: AnalysisListResponse = await apiService.listAnalyses(
                pageSize,
                offset,
                query || undefined
            );
            setAnalyses(response.analyses);
            setTotalCount(response.pagination.total);
            setTotalPages(Math.ceil(response.pagination.total / pageSize));
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to load analyses');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchAnalyses(currentPage, searchQuery);
    }, [currentPage]);

    const handleSearch = () => {
        setCurrentPage(1);
        fetchAnalyses(1, searchQuery);
    };

    const handlePageChange = (_: React.ChangeEvent<unknown>, page: number) => {
        setCurrentPage(page);
    };

    const getScoreColor = (score: number) => {
        if (score >= 80) return theme.palette.success.main;
        if (score >= 60) return theme.palette.warning.main;
        return theme.palette.error.main;
    };

    const getRecommendationChip = (recommendation: string) => {
        const isGo = recommendation.toLowerCase().includes('go') && !recommendation.toLowerCase().includes('no-go');
        return (
            <Chip
                label={isGo ? 'GO' : 'NO-GO'}
                size="small"
                sx={{
                    backgroundColor: isGo ? theme.palette.success.main : theme.palette.error.main,
                    color: 'white',
                    fontWeight: 500,
                }}
            />
        );
    };

    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            <Fade in timeout={800}>
                <Box>
                    {/* Header */}
                    <Box sx={{ mb: 4, textAlign: 'center' }}>
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
                            Analysis Dashboard
                        </Typography>
                        <Typography variant="h6" color="text.secondary" sx={{ mb: 3 }}>
                            View and manage your startup idea analyses
                        </Typography>
                    </Box>

                    {/* Search and Stats */}
                    <Paper
                        elevation={2}
                        sx={{
                            p: 3,
                            mb: 4,
                            background: 'rgba(255, 255, 255, 0.9)',
                            backdropFilter: 'blur(10px)',
                        }}
                    >
                        <Grid container spacing={3} alignItems="center">
                            <Grid item xs={12} md={6}>
                                <TextField
                                    fullWidth
                                    placeholder="Search analyses by title or description..."
                                    value={searchQuery}
                                    onChange={(e) => setSearchQuery(e.target.value)}
                                    onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                                    InputProps={{
                                        startAdornment: (
                                            <InputAdornment position="start">
                                                <Search />
                                            </InputAdornment>
                                        ),
                                        endAdornment: (
                                            <InputAdornment position="end">
                                                <Button
                                                    variant="contained"
                                                    size="small"
                                                    onClick={handleSearch}
                                                    disabled={loading}
                                                >
                                                    Search
                                                </Button>
                                            </InputAdornment>
                                        ),
                                    }}
                                    sx={{
                                        '& .MuiOutlinedInput-root': {
                                            background: 'rgba(255, 255, 255, 0.8)',
                                        },
                                    }}
                                />
                            </Grid>
                            <Grid item xs={12} md={6}>
                                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                    <Typography variant="body1" color="text.secondary">
                                        {loading ? 'Loading...' : `${totalCount} analyses found`}
                                    </Typography>
                                    <Tooltip title="Refresh">
                                        <IconButton onClick={() => fetchAnalyses(currentPage, searchQuery)} disabled={loading}>
                                            <Refresh />
                                        </IconButton>
                                    </Tooltip>
                                </Box>
                            </Grid>
                        </Grid>
                    </Paper>

                    {/* Error State */}
                    {error && (
                        <Alert severity="error" sx={{ mb: 3 }}>
                            {error}
                        </Alert>
                    )}

                    {/* Loading State */}
                    {loading ? (
                        <Box sx={{ textAlign: 'center', py: 8 }}>
                            <CircularProgress size={60} />
                            <Typography variant="h6" sx={{ mt: 2 }}>
                                Loading analyses...
                            </Typography>
                        </Box>
                    ) : (
                        <>
                            {/* Analysis Grid */}
                            {analyses.length === 0 ? (
                                <Box sx={{ textAlign: 'center', py: 8 }}>
                                    <Typography variant="h5" color="text.secondary" sx={{ mb: 2 }}>
                                        No analyses found
                                    </Typography>
                                    <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
                                        {searchQuery ? 'Try adjusting your search terms' : 'Start by submitting your first idea for analysis'}
                                    </Typography>
                                    <Button
                                        variant="contained"
                                        onClick={() => navigate('/')}
                                        sx={{ background: theme.palette.gradient.primary }}
                                    >
                                        Submit New Idea
                                    </Button>
                                </Box>
                            ) : (
                                <>
                                    <Grid container spacing={3} sx={{ mb: 4 }}>
                                        {analyses.map((analysis) => (
                                            <Grid item xs={12} sm={6} lg={4} key={analysis.id}>
                                                <Card
                                                    elevation={2}
                                                    sx={{
                                                        height: '100%',
                                                        background: 'rgba(255, 255, 255, 0.9)',
                                                        backdropFilter: 'blur(10px)',
                                                        border: '1px solid rgba(255, 255, 255, 0.2)',
                                                        transition: 'all 0.3s ease-in-out',
                                                        cursor: 'pointer',
                                                        '&:hover': {
                                                            transform: 'translateY(-4px)',
                                                            boxShadow: theme.shadows[8],
                                                        },
                                                    }}
                                                    onClick={() => navigate(`/analysis/${analysis.id}`)}
                                                >
                                                    <CardContent sx={{ p: 3 }}>
                                                        {/* Header */}
                                                        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
                                                            <Typography
                                                                variant="h6"
                                                                sx={{
                                                                    fontWeight: 600,
                                                                    overflow: 'hidden',
                                                                    textOverflow: 'ellipsis',
                                                                    display: '-webkit-box',
                                                                    WebkitLineClamp: 2,
                                                                    WebkitBoxOrient: 'vertical',
                                                                    flex: 1,
                                                                    mr: 1,
                                                                }}
                                                            >
                                                                {analysis.idea.title}
                                                            </Typography>
                                                            {getRecommendationChip(analysis.verdict.recommendation)}
                                                        </Box>

                                                        {/* Description */}
                                                        <Typography
                                                            variant="body2"
                                                            color="text.secondary"
                                                            sx={{
                                                                mb: 2,
                                                                overflow: 'hidden',
                                                                textOverflow: 'ellipsis',
                                                                display: '-webkit-box',
                                                                WebkitLineClamp: 2,
                                                                WebkitBoxOrient: 'vertical',
                                                                lineHeight: 1.4,
                                                            }}
                                                        >
                                                            {analysis.idea.one_liner}
                                                        </Typography>

                                                        {/* Score */}
                                                        <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                                                            <TrendingUp sx={{ mr: 1, color: getScoreColor(analysis.verdict.overall_score) }} />
                                                            <Typography
                                                                variant="h5"
                                                                sx={{
                                                                    fontWeight: 700,
                                                                    color: getScoreColor(analysis.verdict.overall_score),
                                                                }}
                                                            >
                                                                {analysis.verdict.overall_score.toFixed(1)}
                                                            </Typography>
                                                            <Typography variant="body2" color="text.secondary" sx={{ ml: 0.5 }}>
                                                                /100
                                                            </Typography>
                                                        </Box>

                                                        {/* Tags */}
                                                        <Box sx={{ display: 'flex', gap: 0.5, mb: 2, flexWrap: 'wrap' }}>
                                                            {analysis.idea.category && (
                                                                <Chip label={analysis.idea.category} size="small" variant="outlined" />
                                                            )}
                                                            {analysis.partial && (
                                                                <Chip label="Partial" size="small" color="warning" />
                                                            )}
                                                        </Box>

                                                        {/* Footer */}
                                                        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                                            <Box sx={{ display: 'flex', alignItems: 'center' }}>
                                                                <AccessTime sx={{ fontSize: 16, mr: 0.5, color: 'text.secondary' }} />
                                                                <Typography variant="caption" color="text.secondary">
                                                                    {new Date(analysis.created_at).toLocaleDateString()}
                                                                </Typography>
                                                            </Box>
                                                            <Box sx={{ display: 'flex', gap: 0.5 }}>
                                                                <Tooltip title="View Analysis">
                                                                    <IconButton
                                                                        size="small"
                                                                        onClick={(e) => {
                                                                            e.stopPropagation();
                                                                            navigate(`/analysis/${analysis.id}`);
                                                                        }}
                                                                    >
                                                                        <Visibility fontSize="small" />
                                                                    </IconButton>
                                                                </Tooltip>
                                                            </Box>
                                                        </Box>
                                                    </CardContent>
                                                </Card>
                                            </Grid>
                                        ))}
                                    </Grid>

                                    {/* Pagination */}
                                    {totalPages > 1 && (
                                        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
                                            <Pagination
                                                count={totalPages}
                                                page={currentPage}
                                                onChange={handlePageChange}
                                                color="primary"
                                                size="large"
                                                showFirstButton
                                                showLastButton
                                            />
                                        </Box>
                                    )}
                                </>
                            )}
                        </>
                    )}
                </Box>
            </Fade>
        </Container>
    );
}

