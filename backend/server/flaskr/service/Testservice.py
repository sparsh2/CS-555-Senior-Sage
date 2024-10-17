from unittest.mock import patch, Mock
from service import delete_preferences
import pytest

# @patch('os.remove', Mock())
def test_delete_preferences_raises_error():
    mocked_error = 'insufficient file permission'
    with patch('os.remove', side_effect=OSError(mocked_error)):
        done, err = delete_preferences({'filepath': 'dummy/file/path'})
        assert done == False
        assert err == mocked_error


def test_delete_preferences_deletes_file(mocker):
    mock_os_remove = mocker.patch('os.remove')
    done, err = delete_preferences({'filepath': 'dummy/file/path'})
    mock_os_remove.assert_called()
    assert done == True
    assert err ==""