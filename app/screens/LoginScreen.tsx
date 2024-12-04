import React from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  Button,
  StyleSheet,
  Image,
} from 'react-native';

export default function LoginScreen({ navigation }: any) {
  console.log('Rendering LoginScreen...');
  return (
    <View style={styles.container}>
      {/* Logo and Title */}
      <View style={styles.header}>
        <Image
          source={require('../assets/icons/owl.png')} // Replace with your logo URI
          style={styles.logo}
        />
        <Text style={styles.title}>Senior Sage</Text>
      </View>

      {/* Welcome Text */}
      <Text style={styles.welcomeText}>Welcome</Text>

      {/* Email and Password Fields */}
      <TextInput
        //style={styles.input}
        style={[styles.input, { color: '#000', backgroundColor: '#fff' }]}
        placeholder="Email/Vitalink ID"
        placeholderTextColor="#999"
      />
      <TextInput
        style={styles.input}
        placeholder="Password"
        placeholderTextColor="#999"
        secureTextEntry
      />

      {/* Sign In Button */}
      <TouchableOpacity style={styles.signInButton}>
        <Text style={styles.signInButtonText}>Sign In</Text>
      </TouchableOpacity>

      {/* Forgot Password Link */}
      <TouchableOpacity onPress={() => console.log('Forgot Password pressed')}>
        <Text style={styles.forgotPassword}>Forgot Password?</Text>
      </TouchableOpacity>

      {/* Create Account Button */}
      <TouchableOpacity
        style={styles.createAccountButton}
        onPress={() => navigation.navigate('CreateAccount')} // Navigate to CreateAccount
      >
        <Text style={styles.createAccountButtonText}>Create Account</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
    backgroundColor: '#fff',
  },
  header: {
    alignItems: 'center',
    marginBottom: 20,
  },
  logo: {
    width: 80,
    height: 80,
    marginBottom: 10,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
  },
  welcomeText: {
    fontSize: 28,
    fontWeight: '600',
    marginBottom: 20,
  },
  input: {
    width: '100%',
    height: 50,
    borderColor: '#ccc',
    borderWidth: 1,
    borderRadius: 5,
    paddingHorizontal: 10,
    marginBottom: 15,
  },
  signInButton: {
    width: '100%',
    backgroundColor: '#4CAF50',
    paddingVertical: 15,
    borderRadius: 5,
    alignItems: 'center',
    marginBottom: 10,
  },
  signInButtonText: {
    color: '#fff',
    fontSize: 18,
    fontWeight: 'bold',
  },
  forgotPassword: {
    color: '#007BFF',
    fontSize: 14,
    marginBottom: 20,
    textDecorationLine: 'underline',
  },
  faceId: {
    position: 'absolute',
    top: '55%',
    right: '10%',
  },
  createAccountButton: {
    position: 'absolute',
    bottom: 30,
    borderColor: '#4CAF50',
    borderWidth: 1,
    borderRadius: 5,
    paddingHorizontal: 20,
    paddingVertical: 10,
  },
  createAccountButtonText: {
    color: '#4CAF50',
    fontSize: 16,
    fontWeight: 'bold',
  },
});

