from flask import Flask, render_template, request, redirect, url_for

app = Flask(__name__)

# Fake login credentials
users = {"testuser": "password"}

# Wishlist items
wishlist_items = [
    {
        "brand": "Adidas",
        "name": "Adidas Jacket",
        "price": 99.99,
        "date_added": "2025-01-15",
        "photo_url": "/static/adidas-jacket.jpg",
        "store_url": "https://www.adidas.com",
        "on_sale": False,
    },
    {
        "brand": "Nike",
        "name": "Nike Air Max",
        "price": 120.00,
        "date_added": "2025-01-10",
        "photo_url": "/static/nike-shoes.jpg",
        "store_url": "https://www.nike.com",
        "on_sale": True,
    },
    {
        "brand": "Zara",
        "name": "Zara Blazer",
        "price": 89.99,
        "date_added": "2025-01-05",
        "photo_url": "/static/zara-blazer.jpg",
        "store_url": "https://www.zara.com",
        "on_sale": False,
    },
    {
        "brand": "Zalando",
        "name": "Zalando Exclusive Shirt",
        "price": 75.00,
        "date_added": "2025-01-20",
        "photo_url": "/static/zalando-shirt.jpg",
        "store_url": "https://www.zalando.com",
        "on_sale": False,
    },
    {
        "brand": "Puma",
        "name": "Puma Running Shoes",
        "price": 110.00,
        "date_added": "2025-01-18",
        "photo_url": "/static/puma-shoes.jpg",
        "store_url": "https://www.puma.com",
        "on_sale": False,
    },
    {
        "brand": "H&M",
        "name": "H&M T-Shirt",
        "price": 19.99,
        "date_added": "2025-01-12",
        "photo_url": "/static/hm-tshirt.jpg",
        "store_url": "https://www.hm.com",
        "on_sale": False,
    },
    {
        "brand": "Uniqlo",
        "name": "Uniqlo Chinos",
        "price": 49.99,
        "date_added": "2025-01-08",
        "photo_url": "/static/uniqlo-chinos.jpg",
        "store_url": "https://www.uniqlo.com",
        "on_sale": False,
    },
    {
        "brand": "Levi's",
        "name": "Levi's 501 Jeans",
        "price": 89.50,
        "date_added": "2025-01-07",
        "photo_url": "/static/levis-jeans.jpg",
        "store_url": "https://www.levi.com",
        "on_sale": False,
    },
    {
        "brand": "Gucci",
        "name": "Gucci Belt",
        "price": 450.00,
        "date_added": "2025-01-01",
        "photo_url": "/static/gucci-belt.jpg",
        "store_url": "https://www.gucci.com",
        "on_sale": True,
    },
    {
        "brand": "Prada",
        "name": "Prada Sunglasses",
        "price": 320.00,
        "date_added": "2025-01-03",
        "photo_url": "/static/prada-sunglasses.jpg",
        "store_url": "https://www.prada.com",
        "on_sale": False,
    },
]

# Fake notifications
notifications = [
    {"title": "Item Added", "message": "You added a new item to your wishlist."},
    {"title": "Sale Alert", "message": "Nike Air Max is now on sale!"},
]


@app.route("/")
def home():
    return redirect(url_for("login"))


@app.route("/login", methods=["GET", "POST"])
def login():
    if request.method == "POST":
        username = request.form.get("username")
        password = request.form.get("password")
        if username in users and users[username] == password:
            return redirect(url_for("wishlist"))
        return render_template("login.html", error="Invalid username or password.")
    return render_template("login.html")


@app.route("/wishlist")
def wishlist():
    return render_template("wishlist.html", items=wishlist_items)


@app.route("/notifications")
def notifications_page():
    return render_template("notifications.html", notifications=notifications)


if __name__ == "__main__":
    app.run(debug=True)
