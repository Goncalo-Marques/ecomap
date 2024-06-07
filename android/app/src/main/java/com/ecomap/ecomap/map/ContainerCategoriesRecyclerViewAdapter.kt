package com.ecomap.ecomap.map

import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageView
import android.widget.TextView
import androidx.recyclerview.widget.RecyclerView
import com.ecomap.ecomap.R

/**
 * Represents the structure of an item in the container category recycler view.
 * @param iconResourceID Category icon resource ID.
 * @param category Category description.
 */
data class ContainerCategoryRecyclerViewData(
    val iconResourceID: Int,
    val category: String = "",
)

/**
 * Recycler view adapter for the container categories.
 */
class ContainerCategoriesRecyclerViewAdapter(private val dataSet: Array<ContainerCategoryRecyclerViewData>) :
    RecyclerView.Adapter<ContainerCategoriesRecyclerViewAdapter.ViewHolder>() {

    /**
     * Defines the views in the adapter.
     */
    class ViewHolder(view: View) : RecyclerView.ViewHolder(view) {
        val imageView: ImageView = view.findViewById(R.id.image_view_icon)
        val textView: TextView = view.findViewById(R.id.text_view_category)
    }

    override fun onCreateViewHolder(viewGroup: ViewGroup, viewType: Int): ViewHolder {
        // Create a new view, which defines the UI of the list item.
        val view = LayoutInflater.from(viewGroup.context)
            .inflate(R.layout.container_category, viewGroup, false)

        return ViewHolder(view)
    }

    override fun onBindViewHolder(viewHolder: ViewHolder, position: Int) {
        val data = dataSet[position]

        viewHolder.imageView.setImageResource(data.iconResourceID)
        viewHolder.textView.text = data.category

        if (data.category.isBlank()) {
            // Make the category description invisible if it is not provided.
            viewHolder.textView.visibility = View.GONE
        }
    }

    override fun getItemCount(): Int {
        return dataSet.size
    }
}
