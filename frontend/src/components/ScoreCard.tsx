import {
    Card,
    CardContent,
    Typography,
    Box,
    LinearProgress,
    Chip,
} from '@mui/material';
import { useTheme } from '@mui/material/styles';

interface ScoreCardProps {
    title: string;
    score: number;
    maxScore?: number;
    description?: string;
    level?: string;
    color?: 'primary' | 'secondary' | 'success' | 'warning' | 'error';
}

export default function ScoreCard({
    title,
    score,
    maxScore = 100,
    description,
    level,
    color = 'primary'
}: ScoreCardProps) {
    const theme = useTheme();

    const percentage = (score / maxScore) * 100;

    const getScoreColor = (score: number) => {
        if (score >= 80) return theme.palette.success.main;
        if (score >= 60) return theme.palette.warning.main;
        return theme.palette.error.main;
    };

    const getScoreLevel = (score: number) => {
        if (score >= 80) return 'Excellent';
        if (score >= 60) return 'Good';
        if (score >= 40) return 'Fair';
        return 'Poor';
    };

    const scoreColor = getScoreColor(score);
    const scoreLevel = level || getScoreLevel(score);

    return (
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
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
                    <Typography variant="h6" sx={{ fontWeight: 600, flex: 1 }}>
                        {title}
                    </Typography>
                    <Chip
                        label={scoreLevel}
                        size="small"
                        sx={{
                            backgroundColor: scoreColor,
                            color: 'white',
                            fontWeight: 500,
                        }}
                    />
                </Box>

                <Box sx={{ mb: 2 }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
                        <Typography variant="h3" sx={{ fontWeight: 700, color: scoreColor }}>
                            {score.toFixed(1)}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            /{maxScore}
                        </Typography>
                    </Box>

                    <LinearProgress
                        variant="determinate"
                        value={percentage}
                        sx={{
                            height: 8,
                            borderRadius: 4,
                            backgroundColor: theme.palette.grey[200],
                            '& .MuiLinearProgress-bar': {
                                backgroundColor: scoreColor,
                                borderRadius: 4,
                            },
                        }}
                    />
                </Box>

                {description && (
                    <Typography variant="body2" color="text.secondary" sx={{ lineHeight: 1.5 }}>
                        {description}
                    </Typography>
                )}
            </CardContent>
        </Card>
    );
}

