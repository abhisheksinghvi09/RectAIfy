import React from 'react';
import {
    AppBar,
    Toolbar,
    Typography,
    Button,
    Box,
    Container,
    useScrollTrigger,
    Slide,
} from '@mui/material';
import { useNavigate, useLocation } from 'react-router-dom';
import { useTheme } from '@mui/material/styles';

interface HideOnScrollProps {
    children: React.ReactElement;
}

function HideOnScroll({ children }: HideOnScrollProps) {
    const trigger = useScrollTrigger();

    return (
        <Slide appear={false} direction="down" in={!trigger}>
            {children}
        </Slide>
    );
}

interface LayoutProps {
    children: React.ReactNode;
}

export default function Layout({ children }: LayoutProps) {
    const navigate = useNavigate();
    const location = useLocation();
    const theme = useTheme();

    const isActive = (path: string) => location.pathname === path;

    return (
        <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
            <HideOnScroll>
                <AppBar
                    position="fixed"
                    elevation={0}
                    sx={{
                        background: 'rgba(255, 255, 255, 0.8)',
                        backdropFilter: 'blur(20px)',
                        borderBottom: '1px solid rgba(255, 255, 255, 0.2)',
                        color: theme.palette.text.primary,
                    }}
                >
                    <Container maxWidth="lg">
                        <Toolbar sx={{ px: { xs: 0, sm: 0 } }}>
                            <Typography
                                variant="h6"
                                component="div"
                                sx={{
                                    flexGrow: 1,
                                    fontWeight: 700,
                                    background: theme.palette.gradient.primary,
                                    backgroundClip: 'text',
                                    WebkitBackgroundClip: 'text',
                                    WebkitTextFillColor: 'transparent',
                                    cursor: 'pointer',
                                }}
                                onClick={() => navigate('/')}
                            >
                                RectAIfy
                            </Typography>

                            <Box sx={{ display: 'flex', gap: 2 }}>
                                <Button
                                    color="inherit"
                                    onClick={() => navigate('/')}
                                    sx={{
                                        fontWeight: isActive('/') ? 600 : 400,
                                        color: isActive('/') ? theme.palette.primary.main : 'inherit',
                                        '&:hover': {
                                            background: 'rgba(102, 126, 234, 0.1)',
                                        },
                                    }}
                                >
                                    Submit
                                </Button>
                                <Button
                                    color="inherit"
                                    onClick={() => navigate('/dashboard')}
                                    sx={{
                                        fontWeight: isActive('/dashboard') ? 600 : 400,
                                        color: isActive('/dashboard') ? theme.palette.primary.main : 'inherit',
                                        '&:hover': {
                                            background: 'rgba(102, 126, 234, 0.1)',
                                        },
                                    }}
                                >
                                    Dashboard
                                </Button>
                                <Button
                                    color="inherit"
                                    onClick={() => navigate('/about')}
                                    sx={{
                                        fontWeight: isActive('/about') ? 600 : 400,
                                        color: isActive('/about') ? theme.palette.primary.main : 'inherit',
                                        '&:hover': {
                                            background: 'rgba(102, 126, 234, 0.1)',
                                        },
                                    }}
                                >
                                    About
                                </Button>
                            </Box>
                        </Toolbar>
                    </Container>
                </AppBar>
            </HideOnScroll>

            <Toolbar /> {/* Spacer for fixed AppBar */}

            <Box
                component="main"
                sx={{
                    flexGrow: 1,
                    background: `linear-gradient(135deg, 
            rgba(255, 255, 255, 0.9) 0%, 
            rgba(248, 250, 252, 0.9) 50%, 
            rgba(255, 255, 255, 0.9) 100%
          )`,
                    minHeight: 'calc(100vh - 64px)',
                }}
            >
                {children}
            </Box>
        </Box>
    );
}

