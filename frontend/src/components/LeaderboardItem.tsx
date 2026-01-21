import React from 'react';
import { View, Text, StyleSheet, Pressable } from 'react-native';
import { RankedUser } from '../types';
import { RankBadge } from './RankBadge';
import { colors, spacing, borderRadius, typography } from '../theme';

interface LeaderboardItemProps {
    user: RankedUser;
    onPress?: (user: RankedUser) => void;
}

export const LeaderboardItem: React.FC<LeaderboardItemProps> = ({ user, onPress }) => {
    const isTop3 = user.rank <= 3;

    return (
        <Pressable
            onPress={() => onPress?.(user)}
            style={({ pressed }) => [
                styles.container,
                isTop3 && styles.topRanked,
                pressed && styles.pressed,
            ]}
        >
            <View style={styles.rankContainer}>
                <RankBadge rank={user.rank} />
            </View>

            <View style={styles.userInfo}>
                <Text style={styles.username} numberOfLines={1}>
                    {user.username}
                </Text>
                <Text style={styles.userId}>#{user.id}</Text>
            </View>

            <View style={styles.ratingContainer}>
                <Text style={styles.rating}>{user.rating}</Text>
                <Text style={styles.ratingLabel}>pts</Text>
            </View>
        </Pressable>
    );
};

const styles = StyleSheet.create({
    container: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: colors.card,
        padding: spacing.md,
        marginHorizontal: spacing.md,
        marginVertical: spacing.xs,
        borderRadius: borderRadius.md,
        borderWidth: 1,
        borderColor: colors.border,
    },
    topRanked: {
        borderColor: colors.primaryLight,
        borderWidth: 1.5,
    },
    pressed: {
        opacity: 0.8,
        transform: [{ scale: 0.98 }],
    },
    rankContainer: {
        marginRight: spacing.md,
    },
    userInfo: {
        flex: 1,
    },
    username: {
        color: colors.text,
        fontSize: typography.body.fontSize,
        fontWeight: '600',
    },
    userId: {
        color: colors.textMuted,
        fontSize: typography.small.fontSize,
        marginTop: 2,
    },
    ratingContainer: {
        alignItems: 'flex-end',
    },
    rating: {
        color: colors.primaryLight,
        fontSize: typography.h3.fontSize,
        fontWeight: '700',
    },
    ratingLabel: {
        color: colors.textMuted,
        fontSize: typography.small.fontSize,
    },
});
