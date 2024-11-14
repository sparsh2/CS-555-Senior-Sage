// screens/LoginScreen.tsx
import React from 'react';
import { View, Text, Button } from 'react-native';

export default function LoginScreen({ navigation }: any) {
  const handleLogin = () => {
    // Set login state to true (you can use context or global state management here)
    navigation.replace("Main"); // Navigate to MainScreen after login
  };

  return (
    <View>
      <Text>Login Screen</Text>
      <Button title="Log in" onPress={handleLogin} />
    </View>
  );
}
