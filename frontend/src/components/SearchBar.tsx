import React from 'react';
import { View, TextInput, StyleSheet, Pressable } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { colors, spacing, borderRadius, typography } from '../theme';

interface SearchBarProps {
    value: string;
    onChangeText: (text: string) => void;
    placeholder?: string;
}

export const SearchBar: React.FC<SearchBarProps> = ({
    value,
    onChangeText,
    placeholder = 'Search users...',
}) => {
    return (
        <View style={styles.container}>
            <Ionicons name="search" size={20} color={colors.textMuted} style={styles.icon} />
            <TextInput
                style={styles.input}
                value={value}
                onChangeText={onChangeText}
                placeholder={placeholder}
                placeholderTextColor={colors.textMuted}
                autoCapitalize="none"
                autoCorrect={false}
                returnKeyType="search"
            />
            {value.length > 0 && (
                <Pressable onPress={() => onChangeText('')} style={styles.clearButton}>
                    <Ionicons name="close-circle" size={20} color={colors.textMuted} />
                </Pressable>
            )}
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: colors.card,
        borderRadius: borderRadius.lg,
        paddingHorizontal: spacing.md,
        marginHorizontal: spacing.md,
        marginVertical: spacing.sm,
        borderWidth: 1,
        borderColor: colors.border,
        height: 48,
    },
    icon: {
        marginRight: spacing.sm,
    },
    input: {
        flex: 1,
        color: colors.text,
        fontSize: typography.body.fontSize,
        height: '100%',
    },
    clearButton: {
        padding: spacing.xs,
        marginLeft: spacing.xs,
    },
});
