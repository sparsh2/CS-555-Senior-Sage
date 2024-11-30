import React from 'react';
import { createStackNavigator } from '@react-navigation/stack';
import LoginScreen from '../screens/LoginScreen';
import MainScreen from '../screens/MainScreen';
import CreateAccountScreen from '../screens/CreateAccountScreen';
import ProfilePictureScreen from '../screens/ProfilePictureScreen';
import UserDetailsScreen from '../screens/UserDetailsScreen';

const Stack = createStackNavigator();

export default function AppNavigator() {
  console.log('Initializing AppNavigator...'); // Debugging statement for navigator initialization

  return (
    <Stack.Navigator>
      <Stack.Screen
        name="Login"
        component={LoginScreen}
        options={{ title: 'Login' }}
        listeners={{
          focus: () => console.log('Navigated to Login Screen'), // Debug when entering Login
          blur: () => console.log('Left Login Screen'), // Debug when leaving Login
        }}
      />
      <Stack.Screen
        name="Main"
        component={MainScreen}
        options={{ title: 'Main Screen' }}
        listeners={{
          focus: () => console.log('Navigated to Main Screen'), // Debug when entering Main
          blur: () => console.log('Left Main Screen'), // Debug when leaving Main
        }}
      />
      <Stack.Screen
        name="CreateAccount"
        component={CreateAccountScreen}
        options={{title: 'Create Account'}}
        listeners={{
          focus: () => console.log('Navigate to Create Account Screen'), //Debug when entering Create Account
          blur: () => console.log('Left Create Account Screen'), //Debug when leaving Create Account
        }}
      />
      <Stack.Screen
        name="ProfilePicture"
        component={ProfilePictureScreen}
        options={{title: 'Profile Picture'}}
        listeners={{
          focus: () => console.log('Navigate to Profile Picture Screen'), //Debug when entering Profile Picture Screen
          blur: () => console.log('Left the Profile Picture Screen'), // Debug when leaving Profile Picture Screen
        }}
      />
      <Stack.Screen
        name="UserDetails"
        component={UserDetailsScreen}
        options={{title: 'User Details'}}
        listeners={{
          focus: () => console.log('Navigate to User Details Screen'), //Debug when entering User Details Screen
          blur: () => console.log('Left the UserDetails Screen'), //Debug when leaving the User Details Screen
        }}
      />
    </Stack.Navigator>
  );
}

