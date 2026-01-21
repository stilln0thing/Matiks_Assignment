import React from 'react';
import { NavigationContainer, DefaultTheme } from '@react-navigation/native';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { Ionicons } from '@expo/vector-icons';
import { LeaderboardScreen } from '../screens/LeaderboardScreen';
import { SearchScreen } from '../screens/SearchScreen';
import { colors } from '../theme';

const Tab = createBottomTabNavigator();

const navigationTheme = {
    ...DefaultTheme,
    dark: true,
    colors: {
        ...DefaultTheme.colors,
        primary: colors.primary,
        background: colors.background,
        card: colors.backgroundSecondary,
        text: colors.text,
        border: colors.border,
    },
};

export const AppNavigator: React.FC = () => {
    return (
        <NavigationContainer theme={navigationTheme}>
            <Tab.Navigator
                screenOptions={({ route }) => ({
                    tabBarIcon: ({ focused, color, size }) => {
                        let iconName: keyof typeof Ionicons.glyphMap;

                        if (route.name === 'Leaderboard') {
                            iconName = focused ? 'trophy' : 'trophy-outline';
                        } else if (route.name === 'Search') {
                            iconName = focused ? 'search' : 'search-outline';
                        } else {
                            iconName = 'help-outline';
                        }

                        return <Ionicons name={iconName} size={size} color={color} />;
                    },
                    tabBarActiveTintColor: colors.primary,
                    tabBarInactiveTintColor: colors.textMuted,
                    tabBarStyle: {
                        backgroundColor: colors.backgroundSecondary,
                        borderTopColor: colors.border,
                        paddingBottom: 5,
                        paddingTop: 5,
                        height: 60,
                    },
                    tabBarLabelStyle: {
                        fontSize: 12,
                        fontWeight: '600',
                    },
                    headerShown: false,
                })}
            >
                <Tab.Screen
                    name="Leaderboard"
                    component={LeaderboardScreen}
                    options={{
                        tabBarLabel: 'Leaderboard',
                    }}
                />
                <Tab.Screen
                    name="Search"
                    component={SearchScreen}
                    options={{
                        tabBarLabel: 'Search',
                    }}
                />
            </Tab.Navigator>
        </NavigationContainer>
    );
};
