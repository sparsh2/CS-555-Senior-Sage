import React, { useState } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  ScrollView,
  Alert,
} from 'react-native';

export default function UserDetailsScreen({ navigation }: any) {
  const [name, setName] = useState('');
  const [age, setAge] = useState('');
  const [gender, setGender] = useState('');
  const [medicalConditions, setMedicalConditions] = useState('');
  const [medications, setMedications] = useState('');

  const handleContinue = () => {
    if (!name || !age || !gender) {
      Alert.alert('Validation Error', 'Please fill all required fields.');
      return;
    }

    navigation.navigate('VoiceSelection'); // Replace with your next screen's name
  };

  return (
    <ScrollView contentContainerStyle={styles.container}>
      <Text style={styles.stepText}>Step 2/3</Text>
      <Text style={styles.title}>Help us know you a little better</Text>

      <TextInput
        style={styles.input}
        placeholder="Name"
        placeholderTextColor="#999"
        value={name}
        onChangeText={setName}
      />

      <TextInput
        style={styles.input}
        placeholder="Age (numerical entry only)"
        placeholderTextColor="#999"
        keyboardType="numeric"
        value={age}
        onChangeText={setAge}
      />

      <TouchableOpacity
        style={styles.input}
        onPress={() =>
          Alert.alert('Select Gender', '', [
            { text: 'Male', onPress: () => setGender('Male') },
            { text: 'Female', onPress: () => setGender('Female') },
            { text: 'Other', onPress: () => setGender('Other') },
          ])
        }
      >
        <Text style={gender ? styles.selectedText : styles.placeholderText}>
          {gender || 'Select your gender'}
        </Text>
      </TouchableOpacity>

      <TextInput
        style={styles.input}
        placeholder="Medical Conditions (if any)"
        placeholderTextColor="#999"
        value={medicalConditions}
        onChangeText={setMedicalConditions}
        multiline
        numberOfLines={3}
      />

      <TextInput
        style={styles.input}
        placeholder="Medications (Optional)"
        placeholderTextColor="#999"
        value={medications}
        onChangeText={setMedications}
        multiline
        numberOfLines={3}
      />

      <TouchableOpacity style={styles.continueButton} onPress={handleContinue}>
        <Text style={styles.continueButtonText}>Continue</Text>
      </TouchableOpacity>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flexGrow: 1,
    padding: 20,
    backgroundColor: '#fff',
  },
  stepText: {
    textAlign: 'center',
    fontSize: 14,
    color: '#aaa',
    marginBottom: 10,
  },
  title: {
    fontSize: 20,
    fontWeight: 'bold',
    textAlign: 'center',
    marginBottom: 20,
  },
  input: {
    width: '100%',
    minHeight: 50,
    borderColor: '#ccc',
    borderWidth: 1,
    borderRadius: 5,
    paddingHorizontal: 10,
    marginBottom: 15,
    justifyContent: 'center',
  },
  placeholderText: {
    color: '#999',
  },
  selectedText: {
    color: '#000',
  },
  continueButton: {
    backgroundColor: '#007bff',
    padding: 15,
    borderRadius: 10,
    alignItems: 'center',
    marginTop: 20,
  },
  continueButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: 'bold',
  },
});
