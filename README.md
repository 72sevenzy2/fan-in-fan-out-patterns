# fan-in-fan-out-patterns
a representation of golangs fan-in/fan-out pattern as part of its concurrency. This repo is for my understanding of them.

Fan in -> its when multiple workers send results channels and the results get merged into 1 channel.

Fan out -> when many workers sends many results channels without getting merged.