import React, { useState, useEffect } from 'react';
import { createStackNavigator } from '@react-navigation/stack';
import LoginScreen from '../screens/LoginScreen'; // Adjust paths as necessary
import MainScreen from '../screens/MainScreen';   // Adjust paths as necessary
import { someAsyncLoginCheckFunction } from '../auth';

const Stack = createStackNavigator();

function AppNavigator() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  // Simulate checking user login status
  useEffect(() => {
    const checkLoginStatus = async () => {
      // Replace with your actual login check logic
      const loggedIn = await someAsyncLoginCheckFunction();
      setIsLoggedIn(loggedIn);
    };

    checkLoginStatus();
  }, []);

  return (
    <Stack.Navigator>
      {isLoggedIn ? (
        <Stack.Screen name="Main" component={MainScreen} />
      ) : (
        <Stack.Screen name="Login" component={LoginScreen} />
      )}
    </Stack.Navigator>
  );
}

export default AppNavigator;
