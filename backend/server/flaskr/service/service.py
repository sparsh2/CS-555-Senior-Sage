import os
def delete_preferences(req):
    # TODO: implement this once language model design is finalized
    try:
        os.remove(req['filepath'])
    except OSError as e:
        return False, str(e)
    return True, ""