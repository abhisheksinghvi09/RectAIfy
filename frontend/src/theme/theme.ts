import { createTheme, type ThemeOptions } from '@mui/material/styles';

declare module '@mui/material/styles' {
    interface Palette {
        accent: {
            primary: string;
            secondary: string;
            tertiary: string;
        };
        gradient: {
            primary: string;
            secondary: string;
            tertiary: string;
            surface: string;
        };
    }

    interface PaletteOptions {
        accent?: {
            primary?: string;
            secondary?: string;
            tertiary?: string;
        };
        gradient?: {
            primary?: string;
            secondary?: string;
            tertiary?: string;
            surface?: string;
        };
    }
}

const themeOptions: ThemeOptions = {
    palette: {
        mode: 'light',
        primary: {
            main: '#2563eb', // Modern blue
            light: '#60a5fa',
            dark: '#1d4ed8',
        },
        secondary: {
            main: '#7c3aed', // Elegant purple
            light: '#a78bfa',
            dark: '#5b21b6',
        },
        accent: {
            primary: '#06b6d4', // Cyan accent
            secondary: '#ec4899', // Pink accent  
            tertiary: '#f59e0b', // Amber accent
        },
        gradient: {
            primary: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
            secondary: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
            tertiary: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
            surface: 'linear-gradient(135deg, #ffecd2 0%, #fcb69f 25%, #ffecd2 50%, #fcb69f 75%, #ffecd2 100%)',
        },
        background: {
            default: '#fafafa',
            paper: '#ffffff',
        },
        text: {
            primary: '#1f2937',
            secondary: '#6b7280',
        },
    },
    typography: {
        fontFamily: '"Inter", "SF Pro Display", -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
        h1: {
            fontSize: '3.5rem',
            fontWeight: 700,
            lineHeight: 1.1,
            letterSpacing: '-0.02em',
        },
        h2: {
            fontSize: '2.5rem',
            fontWeight: 600,
            lineHeight: 1.2,
            letterSpacing: '-0.01em',
        },
        h3: {
            fontSize: '2rem',
            fontWeight: 600,
            lineHeight: 1.3,
        },
        h4: {
            fontSize: '1.5rem',
            fontWeight: 500,
            lineHeight: 1.4,
        },
        h5: {
            fontSize: '1.25rem',
            fontWeight: 500,
            lineHeight: 1.5,
        },
        h6: {
            fontSize: '1.125rem',
            fontWeight: 500,
            lineHeight: 1.5,
        },
        body1: {
            fontSize: '1rem',
            lineHeight: 1.6,
        },
        body2: {
            fontSize: '0.875rem',
            lineHeight: 1.5,
        },
        caption: {
            fontSize: '0.75rem',
            fontWeight: 400,
            lineHeight: 1.4,
        },
    },
    shape: {
        borderRadius: 12,
    },
    shadows: [
        'none',
        '0px 2px 4px rgba(0, 0, 0, 0.05)',
        '0px 4px 8px rgba(0, 0, 0, 0.1)',
        '0px 8px 16px rgba(0, 0, 0, 0.1)',
        '0px 12px 24px rgba(0, 0, 0, 0.15)',
        '0px 16px 32px rgba(0, 0, 0, 0.15)',
        '0px 20px 40px rgba(0, 0, 0, 0.2)',
        '0px 24px 48px rgba(0, 0, 0, 0.2)',
        '0px 28px 56px rgba(0, 0, 0, 0.25)',
        '0px 32px 64px rgba(0, 0, 0, 0.25)',
        '0px 36px 72px rgba(0, 0, 0, 0.3)',
        '0px 40px 80px rgba(0, 0, 0, 0.3)',
        '0px 44px 88px rgba(0, 0, 0, 0.35)',
        '0px 48px 96px rgba(0, 0, 0, 0.35)',
        '0px 52px 104px rgba(0, 0, 0, 0.4)',
        '0px 56px 112px rgba(0, 0, 0, 0.4)',
        '0px 60px 120px rgba(0, 0, 0, 0.45)',
        '0px 64px 128px rgba(0, 0, 0, 0.45)',
        '0px 68px 136px rgba(0, 0, 0, 0.5)',
        '0px 72px 144px rgba(0, 0, 0, 0.5)',
        '0px 76px 152px rgba(0, 0, 0, 0.55)',
        '0px 80px 160px rgba(0, 0, 0, 0.55)',
        '0px 84px 168px rgba(0, 0, 0, 0.6)',
        '0px 88px 176px rgba(0, 0, 0, 0.6)',
        '0px 92px 184px rgba(0, 0, 0, 0.65)',
    ],
    components: {
        MuiButton: {
            styleOverrides: {
                root: {
                    textTransform: 'none',
                    borderRadius: 8,
                    padding: '10px 24px',
                    fontSize: '0.95rem',
                    fontWeight: 500,
                    transition: 'all 0.2s ease-in-out',
                    '&:hover': {
                        transform: 'translateY(-1px)',
                        boxShadow: '0px 8px 24px rgba(0, 0, 0, 0.15)',
                    },
                },
                contained: {
                    background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                    '&:hover': {
                        background: 'linear-gradient(135deg, #5a6fd8 0%, #6b4190 100%)',
                    },
                },
            },
        },
        MuiCard: {
            styleOverrides: {
                root: {
                    borderRadius: 16,
                    border: '1px solid rgba(255, 255, 255, 0.1)',
                    backdropFilter: 'blur(20px)',
                    transition: 'all 0.3s ease-in-out',
                    '&:hover': {
                        transform: 'translateY(-4px)',
                        boxShadow: '0px 20px 40px rgba(0, 0, 0, 0.1)',
                    },
                },
            },
        },
        MuiPaper: {
            styleOverrides: {
                root: {
                    backgroundImage: 'none',
                },
                elevation1: {
                    boxShadow: '0px 2px 8px rgba(0, 0, 0, 0.08)',
                },
            },
        },
        MuiTextField: {
            styleOverrides: {
                root: {
                    '& .MuiOutlinedInput-root': {
                        borderRadius: 12,
                        transition: 'all 0.2s ease-in-out',
                        '&:hover': {
                            '& .MuiOutlinedInput-notchedOutline': {
                                borderColor: '#667eea',
                            },
                        },
                        '&.Mui-focused': {
                            '& .MuiOutlinedInput-notchedOutline': {
                                borderColor: '#667eea',
                                borderWidth: 2,
                            },
                        },
                    },
                },
            },
        },
        MuiLinearProgress: {
            styleOverrides: {
                root: {
                    borderRadius: 4,
                    height: 8,
                },
            },
        },
    },
};

export const theme = createTheme(themeOptions);
export default theme;

