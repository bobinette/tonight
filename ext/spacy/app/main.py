#!/usr/bin/env
# -*-coding:Utf-8 -*

# =============================================================================
# IMPORTS
# =============================================================================

from flask import Flask, abort, request, jsonify

import spacy


# =============================================================================
# Routes
# =============================================================================

app = Flask(__name__)

# Use the basic en model. If you want to use another one, change
# the name here.
# If you use docker, make sure to update the Dockerifle accordingly
nlp = spacy.load('en')

debug = True


@app.route('/similarities', methods=['POST'])
def similarities():
    if not request.json:
        abort(400, 'no json body')

    tasks = request.json.get('tasks')
    if not tasks:
        abort(400, 'missing tasks in body')

    sims = []
    for i, t1 in enumerate(tasks[:-1]):
        for t2 in tasks[i + 1:]:
            doc1 = nlp(t1)
            doc2 = nlp(t2)
            sim = doc1.similarity(doc2)

            sims.append({
                'left': t1,
                'right': t2,
                'similarity': sim
            })

    return jsonify(sims)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=1717, debug=debug, threaded=True)
