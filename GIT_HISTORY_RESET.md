# Git History Reset Instructions - Based on Current State

## Current State Analysis
✅ Remote is already set to: `https://github.com/ozacod/cpx.git`
✅ You're currently on `main` branch
✅ Working tree is clean
✅ Already have "Initial commit - v1.0.0" (de1069d)
⚠️ Both `master` and `main` branches exist locally and remotely
⚠️ **master is still the default branch on GitHub** (this is why deletion failed)

## Instructions to Complete the Reset

### Step 1: Change Default Branch on GitHub (MUST DO FIRST)

You **must** change the default branch on GitHub before you can delete master:

1. Go to: https://github.com/ozacod/cpx/settings/branches
2. Under "Default branch" section, you'll see `master` is currently selected
3. Click the switch/change button next to `master`
4. Select `main` from the dropdown
5. Click "Update"
6. Confirm the change (GitHub will warn you about changing the default branch)

**Alternative method:**
- Go to: https://github.com/ozacod/cpx/branches
- Find the `main` branch
- Click the gear icon (⚙️) next to it
- Select "Set as default branch"

### Step 2: Delete the master branch locally
```bash
git branch -D master
```

### Step 3: Delete the master branch on remote (NOW it will work)
```bash
git push origin --delete master
```

### Step 4: Verify you're on main and it's up to date
```bash
git checkout main
git pull origin main
```

### Step 5: If you want to ensure main has only the v1.0.0 commit (fresh start)
If the remote still has old history, you can force push your clean main:
```bash
# This will overwrite remote main with your local clean version
git push -f origin main
```

### Step 6: Tag the v1.0.0 release (if not already done)
```bash
git tag v1.0.0
git push origin v1.0.0
```

### Step 7: Verify everything is clean
```bash
# Check branches
git branch -a
# Should only show main (no master)

# Check remote
git remote -v
# Should show: origin  https://github.com/ozacod/cpx.git

# Check tags
git tag
# Should show: v1.0.0

# Check commit history
git log --oneline
# Should show only: de1069d Initial commit - v1.0.0
```

## If You Need to Start Completely Fresh

If the remote still has old commits and you want to completely reset:

```bash
# 1. Create orphan branch (fresh start)
git checkout --orphan fresh-main

# 2. Add all current files
git add -A

# 3. Commit as v1.0.0
git commit -m "Initial commit - v1.0.0"

# 4. Delete old main locally
git branch -D main

# 5. Rename to main
git branch -m main

# 6. Force push (WARNING: Deletes all remote history)
git push -f origin main

# 7. Change default branch on GitHub to main (see Step 1 above)

# 8. Delete master on remote
git push origin --delete master

# 9. Tag the release
git tag v1.0.0
git push origin v1.0.0
```

## Quick Commands Summary (After Changing Default Branch)

```bash
# 1. FIRST: Change default branch on GitHub (use web interface)
# Go to: https://github.com/ozacod/cpx/settings/branches

# 2. Then delete master branches
git branch -D master
git push origin --delete master

# 3. Ensure main is clean and tagged
git checkout main
git tag v1.0.0
git push origin v1.0.0 --tags

# 4. Verify
git branch -a
git log --oneline
git tag
```

## Important Notes

⚠️ **WARNING**: 
- Force pushing will delete all commit history on the remote
- Deleting branches is permanent
- Anyone who cloned the repo will need to re-clone or reset their local copy
- **You cannot delete the default branch** - must change default first!

✅ **Safe to do**:
- You're already on main branch
- Remote is already pointing to cpx repository
- Working tree is clean

## Troubleshooting

**Error: "refusing to delete the current branch"**
- This means master is still the default branch on GitHub
- Solution: Change default branch to main first (Step 1 above)
- Then you can delete master
