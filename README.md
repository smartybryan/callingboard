# CALLORG

This project is the beginnings of a tool to
assist Ward Bishopric with re-organizing their ward.

It will be a way to model proposed changes
and track releasings and sustainings.

TODO
- Add sustain list on left
- Member box to drag a member to a calling; add to sustain list.
- Style existing calling box same as the callings in the tree
- In drag and drop operations, call the backend API.
- Save and load model.
- Reload original data.
- Refresh tree with new model (necessary?)
- Search box for members (needed?)
- Authentication/authorization ?

WORKFLOW
- Release and Sustain box
- Drag a calling into Release to release the member (make the calling vacant)
- Drag a member (from the member box) to a calling in the tree and it adds them to the Sustain box
- Delete a row from either box to back out the change.
** The drag and drop process will need to call the backend Api to keep the model in sync with the UI.
We won't be able to repaint the tree or the expand state will be lost.
- 

