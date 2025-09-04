import React, { useState } from 'react';
import {
    Container,
    Typography,
    Box,
    Card,
    CardContent,
    TextField,
    Button,
    Grid,
    Chip,
    Alert,
    CircularProgress,
    Fade,
    Slide,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { useTheme } from '@mui/material/styles';
import { apiService } from '../services/apiService';
import { type AnalysisRequest } from '../types/api';

export default function HomePage() {
    const [formData, setFormData] = useState({
        title: '',
        one_liner: '',
        category: '',
        location: '',
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();
    const theme = useTheme();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setLoading(true);

        try {
            const request: AnalysisRequest = {
                idea: {
                    title: formData.title.trim(),
                    one_liner: formData.one_liner.trim(),
                    category: formData.category.trim() || undefined,
                    location: formData.location.trim() || undefined,
                },
            };

            const response = await apiService.submitAnalysis(request);
            navigate(`/analysis/${response.analysis_id}`);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to submit analysis');
        } finally {
            setLoading(false);
        }
    };

    const handleChange = (field: string) => (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData(prev => ({ ...prev, [field]: e.target.value }));
    };

    const isFormValid = formData.title.trim() && formData.one_liner.trim().length >= 10;

    const categories = [
        'SaaS', 'E-commerce', 'FinTech', 'HealthTech', 'EdTech',
        'AI/ML', 'Blockchain', 'IoT', 'Social', 'Gaming', 'Other'
    ];

    return (
        <Container maxWidth="md" sx={{ py: { xs: 4, md: 8 } }}>
            <Fade in timeout={800}>
                <Box>
                    {/* Hero Section */}
                    <Box sx={{ textAlign: 'center', mb: 6 }}>
                        <Slide direction="down" in timeout={1000}>
                            <Typography
                                variant="h1"
                                sx={{
                                    background: theme.palette.gradient.primary,
                                    backgroundClip: 'text',
                                    WebkitBackgroundClip: 'text',
                                    WebkitTextFillColor: 'transparent',
                                    mb: 2,
                                    fontSize: { xs: '2.5rem', md: '3.5rem' },
                                }}
                            >
                                Validate Your Startup Idea
                            </Typography>
                        </Slide>

                        <Slide direction="up" in timeout={1200}>
                            <Typography
                                variant="h5"
                                sx={{
                                    color: theme.palette.text.secondary,
                                    mb: 4,
                                    maxWidth: '600px',
                                    mx: 'auto',
                                    lineHeight: 1.6,
                                }}
                            >
                                Get comprehensive AI-powered analysis across market, problem,
                                barriers, execution, risks, and competitive landscape.
                            </Typography>
                        </Slide>
                    </Box>

                    {/* Submission Form */}
                    <Slide direction="up" in timeout={1400}>
                        <Card
                            elevation={3}
                            sx={{
                                background: 'rgba(255, 255, 255, 0.9)',
                                backdropFilter: 'blur(20px)',
                                border: '1px solid rgba(255, 255, 255, 0.2)',
                                borderRadius: 3,
                            }}
                        >
                            <CardContent sx={{ p: { xs: 3, md: 4 } }}>
                                <Typography
                                    variant="h4"
                                    sx={{
                                        mb: 3,
                                        fontWeight: 600,
                                        color: theme.palette.text.primary,
                                    }}
                                >
                                    Submit Your Idea
                                </Typography>

                                {error && (
                                    <Alert severity="error" sx={{ mb: 3 }}>
                                        {error}
                                    </Alert>
                                )}

                                <Box component="form" onSubmit={handleSubmit}>
                                    <Grid container spacing={3}>
                                        <Grid item xs={12}>
                                            <TextField
                                                fullWidth
                                                label="Idea Title"
                                                placeholder="e.g., TaskAI"
                                                value={formData.title}
                                                onChange={handleChange('title')}
                                                required
                                                variant="outlined"
                                                sx={{
                                                    '& .MuiOutlinedInput-root': {
                                                        background: 'rgba(255, 255, 255, 0.8)',
                                                    },
                                                }}
                                            />
                                        </Grid>

                                        <Grid item xs={12}>
                                            <TextField
                                                fullWidth
                                                label="One-liner Description"
                                                placeholder="e.g., AI-powered task automation platform for small businesses"
                                                value={formData.one_liner}
                                                onChange={handleChange('one_liner')}
                                                required
                                                multiline
                                                rows={3}
                                                variant="outlined"
                                                helperText="Minimum 10 characters. Be specific about what your idea does."
                                                sx={{
                                                    '& .MuiOutlinedInput-root': {
                                                        background: 'rgba(255, 255, 255, 0.8)',
                                                    },
                                                }}
                                            />
                                        </Grid>

                                        <Grid item xs={12} md={6}>
                                            <TextField
                                                fullWidth
                                                label="Category (Optional)"
                                                placeholder="Select or type category"
                                                value={formData.category}
                                                onChange={handleChange('category')}
                                                variant="outlined"
                                                sx={{
                                                    '& .MuiOutlinedInput-root': {
                                                        background: 'rgba(255, 255, 255, 0.8)',
                                                    },
                                                }}
                                            />
                                            <Box sx={{ mt: 1, display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                                                {categories.map((category) => (
                                                    <Chip
                                                        key={category}
                                                        label={category}
                                                        size="small"
                                                        onClick={() => setFormData(prev => ({ ...prev, category }))}
                                                        sx={{
                                                            cursor: 'pointer',
                                                            '&:hover': {
                                                                background: theme.palette.primary.light,
                                                                color: 'white',
                                                            },
                                                        }}
                                                    />
                                                ))}
                                            </Box>
                                        </Grid>

                                        <Grid item xs={12} md={6}>
                                            <TextField
                                                fullWidth
                                                label="Location (Optional)"
                                                placeholder="e.g., San Francisco, US"
                                                value={formData.location}
                                                onChange={handleChange('location')}
                                                variant="outlined"
                                                sx={{
                                                    '& .MuiOutlinedInput-root': {
                                                        background: 'rgba(255, 255, 255, 0.8)',
                                                    },
                                                }}
                                            />
                                        </Grid>

                                        <Grid item xs={12}>
                                            <Button
                                                type="submit"
                                                variant="contained"
                                                size="large"
                                                disabled={!isFormValid || loading}
                                                sx={
                                                    {
                                                        py: 1.5,
                                                        px: 4,
                                                        background: theme.palette.gradient.primary,
                                                        fontSize: '1.1rem',
                                                        fontWeight: 600,
                                                        '&:hover': {
                                                            background: theme.palette.gradient.primary,
                                                            transform: 'translateY(-2px)',
                                                        },
                                                        '&:disabled': {
                                                            background: theme.palette.action.disabledBackground,
                                                            color: theme.palette.action.disabled,
                                                        },
                                                    }
                                                }
                                            >
                                                {loading ? (
                                                    <>
                                                        <CircularProgress size={24} sx={{ mr: 2 }} />
                                                        Analyzing...
                                                    </>
                                                ) : (
                                                    'Analyze Idea'
                                                )}
                                            </Button>
                                        </Grid>
                                    </Grid>
                                </Box>
                            </CardContent>
                        </Card>
                    </Slide>

                    {/* Features Preview */}
                    <Slide direction="up" in timeout={1600}>
                        <Box sx={{ mt: 8, textAlign: 'center' }}>
                            <Typography variant="h5" sx={{ mb: 4, fontWeight: 500 }}>
                                What You'll Get
                            </Typography>
                            <Grid container spacing={3}>
                                {[
                                    { title: 'Market Analysis', desc: 'Competition and positioning insights' },
                                    { title: 'Problem Validation', desc: 'Evidence-based pain point assessment' },
                                    { title: 'Execution Analysis', desc: 'Resource and complexity evaluation' },
                                    { title: 'Risk Assessment', desc: 'Potential challenges and mitigation' },
                                    { title: 'Graveyard Study', desc: 'Learn from similar failed ventures' },
                                    { title: 'Viability Score', desc: 'Overall recommendation with insights' },
                                ].map((feature, index) => (
                                    <Grid item xs={12} sm={6} md={4} key={index}>
                                        <Card
                                            elevation={1}
                                            sx={{
                                                p: 3,
                                                height: '100%',
                                                background: 'rgba(255, 255, 255, 0.6)',
                                                '&:hover': {
                                                    transform: 'translateY(-4px)',
                                                    boxShadow: theme.shadows[4],
                                                },
                                            }}
                                        >
                                            <Typography variant="h6" sx={{ mb: 1, fontWeight: 600 }}>
                                                {feature.title}
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                {feature.desc}
                                            </Typography>
                                        </Card>
                                    </Grid>
                                ))}
                            </Grid>
                        </Box>
                    </Slide>
                </Box>
            </Fade>
        </Container>
    );
}

