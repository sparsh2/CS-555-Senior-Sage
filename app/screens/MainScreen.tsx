// screens/MainScreen.tsx
import React from 'react';
import { View, Text, Button } from 'react-native';

export default function MainScreen({ navigation }: any) {
  const handleLogout = () => {
    navigation.replace("Login"); // Navigate back to LoginScreen on logout
  };

  return (
    <View>
      <Text>Main Screen</Text>
      <Button title="Log out" onPress={handleLogout} />
    </View>
  );
}
