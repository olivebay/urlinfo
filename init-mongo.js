db.blacklists.drop();
db.blacklists.insertOne(
    {
        domain: 'domain.com',
        blacklists: ["vault", "spamHaus"],
    }
)