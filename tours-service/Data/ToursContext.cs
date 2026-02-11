using Microsoft.EntityFrameworkCore;
using TourService.Domain;

namespace TourService.Data;

public class ToursContext : DbContext
{
    public ToursContext(DbContextOptions<ToursContext> options) : base(options) {}

    public DbSet<Tour> Tours => Set<Tour>();
    public DbSet<Tag> Tags => Set<Tag>();
    public DbSet<TourTag> TourTags => Set<TourTag>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.Entity<Tour>(e =>
        {
            e.Property(x => x.Price).HasPrecision(10, 2);
            e.Property(x => x.CreatedAtUtc).HasConversion(v => v, v => DateTime.SpecifyKind(v, DateTimeKind.Utc));
            e.Property(x => x.UpdatedAtUtc).HasConversion(v => v, v => DateTime.SpecifyKind(v, DateTimeKind.Utc));
        });

        modelBuilder.Entity<Tag>(e =>
        {
            e.HasIndex(x => x.Name).IsUnique();
        });

        modelBuilder.Entity<TourTag>(e =>
        {
            e.HasKey(x => new { x.TourId, x.TagId });

            e.HasOne(x => x.Tour)
                .WithMany(t => t.TourTags)
                .HasForeignKey(x => x.TourId);

            e.HasOne(x => x.Tag)
                .WithMany(t => t.TourTags)
                .HasForeignKey(x => x.TagId);
        });
    }
}
