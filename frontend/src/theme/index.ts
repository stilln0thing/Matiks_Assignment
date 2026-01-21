// Theme constants for consistent styling

export const colors = {
    // Primary palette
    primary: '#6366f1',
    primaryDark: '#4f46e5',
    primaryLight: '#818cf8',

    // Background
    background: '#0f0f23',
    backgroundSecondary: '#1a1a2e',
    card: '#16213e',

    // Text
    text: '#ffffff',
    textSecondary: '#94a3b8',
    textMuted: '#64748b',

    // Accents
    gold: '#fbbf24',
    silver: '#9ca3af',
    bronze: '#cd7f32',

    // Status
    success: '#10b981',
    error: '#ef4444',
    warning: '#f59e0b',

    // Border
    border: '#1e293b',
    borderLight: '#334155',
};

export const spacing = {
    xs: 4,
    sm: 8,
    md: 16,
    lg: 24,
    xl: 32,
};

export const typography = {
    h1: {
        fontSize: 28,
        fontWeight: '700' as const,
    },
    h2: {
        fontSize: 22,
        fontWeight: '600' as const,
    },
    h3: {
        fontSize: 18,
        fontWeight: '600' as const,
    },
    body: {
        fontSize: 16,
        fontWeight: '400' as const,
    },
    caption: {
        fontSize: 14,
        fontWeight: '400' as const,
    },
    small: {
        fontSize: 12,
        fontWeight: '400' as const,
    },
};

export const borderRadius = {
    sm: 6,
    md: 12,
    lg: 16,
    xl: 24,
    full: 9999,
};
