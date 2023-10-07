import argparse
import numpy as np
from keras.models import load_model
from keras.preprocessing.sequence import pad_sequences
import pickle
import pandas as pd

# Set the absolute paths to the model and tokenizer files
model_path = '../ml-model/trained_nn_model.h5'
tokenizer_path = '../ml-model/tokenizer.pickle'

# Load the trained model and tokenizer
model = load_model(model_path)
with open(tokenizer_path, 'rb') as handle:
    tokenizer = pickle.load(handle)

# Define your labels
labels = {
    0: 'Sadness',
    1: 'Joy',
    2: 'Love',
    3: 'Anger',
    4: 'Fear',
    5: 'Surprise'
}

def predict_emotions(csv_file):
    # Read text data from the CSV file
    data = pd.read_csv(csv_file)
    text = data['text']

    # Tokenize and preprocess the input text
    sequences = tokenizer.texts_to_sequences(text)
    input_data = pad_sequences(sequences, maxlen=100)

    # Predict emotions
    predictions = model.predict(input_data)
    predicted_labels = [labels[np.argmax(prediction)] for prediction in predictions]

    # Print the predicted emotions for each text
    # for i, text_entry in enumerate(text):
    #     print(f'Text: {text_entry}')
    #     print(f'Predicted Emotion: {predicted_labels[i]}\n')

    # Calculate the average emotion scores
    emotion_counts = {label: 0 for label in labels.values()}
    for label in predicted_labels:
        emotion_counts[label] += 1

    total_entries = len(text)
    average_emotions = {label: count / total_entries for label, count in emotion_counts.items()}

    # Return the average emotions as a formatted string
    result_string = ""
    for label, score in average_emotions.items():
        result_string += f'{label}: {score:.2f}\n'
    
    return result_string
     

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Predict emotions from a CSV file and display the average emotions.")
    parser.add_argument("csv_file", help="Path to the CSV file containing text data")
    args = parser.parse_args()

    try:
        emotions_result = predict_emotions(args.csv_file)
        print(emotions_result)
    except Exception as e:
        print(f"Error: {str(e)}")
