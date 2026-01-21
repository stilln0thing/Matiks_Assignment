import React from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { colors, spacing, borderRadius, typography } from '../theme';

interface RankBadgeProps {
    rank: number;
    size?: 'small' | 'medium' | 'large';
}

export const RankBadge: React.FC<RankBadgeProps> = ({ rank, size = 'medium' }) => {
    const getBackgroundColor = () => {
        if (rank === 1) return colors.gold;
        if (rank === 2) return colors.silver;
        if (rank === 3) return colors.bronze;
        return colors.backgroundSecondary;
    };

    const getTextColor = () => {
        if (rank <= 3) return '#000000';
        return colors.textSecondary;
    };

    const getSize = () => {
        switch (size) {
            case 'small': return { width: 28, height: 28, fontSize: 12 };
            case 'large': return { width: 48, height: 48, fontSize: 20 };
            default: return { width: 36, height: 36, fontSize: 14 };
        }
    };

    const sizeStyles = getSize();

    // Show trophy for top 3
    if (rank <= 3) {
        return (
            <View style={[
                styles.badge,
                {
                    backgroundColor: getBackgroundColor(),
                    width: sizeStyles.width,
                    height: sizeStyles.height,
                }
            ]}>
                <Ionicons name="trophy" size={sizeStyles.fontSize + 2} color="#000" />
            </View>
        );
    }

    return (
        <View style={[
            styles.badge,
            {
                backgroundColor: getBackgroundColor(),
                width: sizeStyles.width,
                height: sizeStyles.height,
            }
        ]}>
            <Text style={[styles.text, { color: getTextColor(), fontSize: sizeStyles.fontSize }]}>
                {rank}
            </Text>
        </View>
    );
};

const styles = StyleSheet.create({
    badge: {
        borderRadius: borderRadius.full,
        alignItems: 'center',
        justifyContent: 'center',
    },
    text: {
        fontWeight: '700',
    },
});
