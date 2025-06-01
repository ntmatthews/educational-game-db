#!/bin/bash

# Educational Game Database - Remote Repository Setup Script
# This script helps set up remote repositories for GitHub or GitLab

echo "🎓 Educational Game Database - Remote Setup"
echo "==========================================="
echo ""

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo "❌ Error: Not in a git repository. Please run this from the project root."
    exit 1
fi

echo "Choose your remote repository provider:"
echo "1) GitHub"
echo "2) GitLab"
echo "3) Both"
echo ""
read -p "Enter your choice (1-3): " choice

case $choice in
    1)
        echo ""
        echo "📝 GitHub Setup Instructions:"
        echo "1. Go to https://github.com/new"
        echo "2. Repository name: educational-game-db"
        echo "3. Description: Educational game student database system with CLI and web interface"
        echo "4. Choose Public or Private"
        echo "5. Don't initialize with README (we already have one)"
        echo ""
        read -p "Enter your GitHub username: " github_user
        read -p "Repository created? Press Enter to continue..."
        
        git remote add origin "https://github.com/$github_user/educational-game-db.git"
        echo "🔗 Added GitHub remote"
        ;;
    2)
        echo ""
        echo "📝 GitLab Setup Instructions:"
        echo "1. Go to https://gitlab.com/projects/new"
        echo "2. Project name: educational-game-db"
        echo "3. Description: Educational game student database system with CLI and web interface"
        echo "4. Choose Public or Private"
        echo "5. Don't initialize with README (we already have one)"
        echo ""
        read -p "Enter your GitLab username: " gitlab_user
        read -p "Repository created? Press Enter to continue..."
        
        git remote add gitlab "https://gitlab.com/$gitlab_user/educational-game-db.git"
        echo "🔗 Added GitLab remote"
        ;;
    3)
        echo ""
        echo "📝 Setting up both GitHub and GitLab:"
        echo ""
        echo "GitHub Setup:"
        echo "1. Go to https://github.com/new"
        echo "2. Repository name: educational-game-db"
        echo "3. Don't initialize with README"
        echo ""
        read -p "Enter your GitHub username: " github_user
        read -p "GitHub repository created? Press Enter to continue..."
        
        echo ""
        echo "GitLab Setup:"
        echo "1. Go to https://gitlab.com/projects/new"
        echo "2. Project name: educational-game-db"
        echo "3. Don't initialize with README"
        echo ""
        read -p "Enter your GitLab username: " gitlab_user
        read -p "GitLab repository created? Press Enter to continue..."
        
        git remote add origin "https://github.com/$github_user/educational-game-db.git"
        git remote add gitlab "https://gitlab.com/$gitlab_user/educational-game-db.git"
        echo "🔗 Added both GitHub and GitLab remotes"
        ;;
    *)
        echo "❌ Invalid choice. Exiting."
        exit 1
        ;;
esac

echo ""
echo "🚀 Pushing to remote repository/repositories..."

# Push to GitHub if origin exists
if git remote get-url origin &> /dev/null; then
    echo "📤 Pushing to GitHub..."
    git push -u origin main
    echo "✅ Pushed to GitHub successfully!"
fi

# Push to GitLab if gitlab remote exists
if git remote get-url gitlab &> /dev/null; then
    echo "📤 Pushing to GitLab..."
    git push -u gitlab main
    echo "✅ Pushed to GitLab successfully!"
fi

echo ""
echo "🎉 Remote repository setup complete!"
echo ""
echo "📋 Next steps:"
echo "- Your code is now backed up to remote repository/repositories"
echo "- You can clone this repository on other machines"
echo "- Set up branch protection rules if needed"
echo "- Consider setting up CI/CD pipelines"
echo ""
echo "🔗 Remote URLs:"
git remote -v
